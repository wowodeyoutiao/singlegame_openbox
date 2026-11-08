using Cysharp.Threading.Tasks;
using DG.Tweening;
using HotUpdateScripts;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using Spine.Unity;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using static GameEvent;

namespace HallBattle
{
    public enum HallRoleState
    {
        None = 0,
        Rest, // 休息
        RestToBattle, // 休息室去战场的转场
        MoveToBattlefield, // 跑向战场
        Waiting, // 等待怪物
        Attack, // 攻击
        TakingDamage, // 受击
        Revive, // 复活
        WaitingDie, // 等待死亡
        Die, // 死亡
        WaitingSettlement, // 结算等待
    }

    public class HallBaseRole : MonoBehaviour
    {
        protected class ShowSkillInfo
        {
            public CFGItem_SkillMainBattle skillCfg;
            private float releaseTime;
            private float timeScale;

            public ShowSkillInfo(CFGItem_SkillMainBattle cfg, float tScale)
            {
                skillCfg = cfg;
                timeScale = tScale;
            }

            public void Act()
            {
                releaseTime = Time.time;
            }

            public bool CanRelease()
            {
                return Time.time >= releaseTime + skillCfg.CD * 0.001f / timeScale;
            }
        }

        [Header("基本属性")]
        protected bool isMainRole; // 主角
        protected int configId; // 角色表ID
        protected long uid;
        protected int soltId;
        protected int job;
        protected float scale;

        protected float maxHealth = 100; // 初始血量
        protected float currentHealth; // 当前血量
        protected float limitX; // 战斗位置，最多可移动的地方

        // 移速、攻击范围
        protected float moveSpeed = 1f;
        protected float attackRange = 300f; // 实际长度 * 0.425，角色的缩放

        protected float timeScale = 1f; // 攻速
        protected float attackCooldown = 1f;// 攻击冷却
        protected float lastAttackTime; // 上次攻击时间

        private Vector2 moveDir = new Vector2(1, 0);

        [Header("战斗状态")]
        public HallRoleState curRoleState; // 角色状态

        public Rigidbody2D rb;
        protected BoxCollider2D boxCollider;

        protected string curSpineName; // spine动画名字，防止重复
        protected Transform spineTrans; // spine动画
        protected List<HallMonster> targets = new List<HallMonster>(); // 攻击目标

        protected List<ShowSkillInfo> norRoleSkillInfos = new List<ShowSkillInfo>(); // 普通技能，非特技

        // 飞血容器
        private GameObject hpChangeGo;

        // 攻击次数，给主线用的
        protected int attackTimes;
        protected int specialSkillCfgID;

        protected int attackRound;

        // 暂停
        protected bool isPause;

        protected virtual void Awake()
        {
            LoadRB();
            UniGameEvent.AddEventListener<GameEvent.HallBallPauseOrPlay>(HallBallPauseOrPlayEvt);
        }

        protected virtual void OnDestroy()
        {
            UniGameEvent.RemoveEventListener<GameEvent.HallBallPauseOrPlay>(HallBallPauseOrPlayEvt);
        }

        private void HallBallPauseOrPlayEvt(HallBallPauseOrPlay arge)
        {
            isPause = arge.pause;
        }

        private void LoadRB()
        {
            if (rb == null)
            {
                rb = gameObject.AddComponent<Rigidbody2D>();
                rb.bodyType = RigidbodyType2D.Kinematic;

                boxCollider = gameObject.AddComponent<BoxCollider2D>();
                boxCollider.isTrigger = true;

                boxCollider.offset = new Vector2(-80, 80);
                boxCollider.size = new Vector2(160, 300);
            }
        }

        #region 初始化数据、位置等
        /// <summary>
        /// 初始化信息
        /// </summary>
        /// <param name="birthPosition"></param>
        public virtual void InitWithInfo(bool isMainRole, int job, int configId = 0, int uid = 0, int soltId = 0, float tScale = 1)
        {
            this.isMainRole = isMainRole;
            this.job = job;
            this.configId = configId;
            this.uid = uid;
            this.soltId = soltId;
            timeScale = tScale;
            scale = transform.localScale.x;

            limitX = -140;
            if (job != 1)
            {
                limitX = -200;
                attackRange = 300 * 1.3f;
            }

            PreLoadRoleSkills();
        }

        public virtual void InitData()
        {
            if (isMainRole)
                maxHealth = RoleM.Instance.GetRoleAttrValueByAttr(pb.Attribute.MAX_HEALTH);
            else
                maxHealth = PartnerM.Instance.GetRoleAttrValueByAttr(uid, pb.Attribute.MAX_HEALTH);

            currentHealth = maxHealth;

            curSpineName = string.Empty;
            targets.Clear();
            attackRound = 1;

            isPause = false;

            LoadRB();
        }

        /// <summary>
        ///  设置spine transform
        /// </summary>
        public virtual void SetSpineTransform(Transform trans)
        {
            spineTrans = trans;
        }

        /// <summary>
        /// 初始化位置等
        /// </summary>
        /// <param name="birthPosition"></param>
        public virtual void InitPosition(Vector2 birthPosition)
        {
            transform.localPosition = birthPosition;
        }

        public virtual void AddHpChangeTrans(Transform trans)
        {
            hpChangeGo = trans.GetChild(2).gameObject;
        }

        /// <summary>
        /// 预加载技能
        /// </summary>
        public virtual void PreLoadRoleSkills()
        {
            norRoleSkillInfos.Clear();
            var skills = CFG_SkillMainBattle.GetNormalSkillsByRoleVocation(RoleM.Instance.Job, RoleM.Instance.Level);
            if (skills != null && skills.Count > 0)
            {
                for (int i = 0; i < skills.Count; i++)
                {
                    norRoleSkillInfos.Add(new ShowSkillInfo(skills[i], timeScale));
                }
            }
        }
        #endregion

        #region 战场战斗
        // 添加目标
        public virtual void AddTarget(HallMonster monster)
        {
            targets.Add(monster);
        }

        protected virtual void Update()
        {
            if (isPause) return;

            if (CheckWin())
                return;

            if (curRoleState >= HallRoleState.MoveToBattlefield && curRoleState <= HallRoleState.Attack)
            {
                var target = GetNearestTarget();
                if (curRoleState == HallRoleState.MoveToBattlefield)
                {
                    if (transform.localPosition.x >= limitX)
                        ChangeRoleState(HallRoleState.Waiting);
                    
                    if (target != null && IsInAttackRange(target))
                    {
                        ChangeRoleState(HallRoleState.Attack);
                        TryAttack(target);
                    }
                }
                else if (curRoleState == HallRoleState.Waiting)
                {
                    if (transform.localPosition.x >= limitX)
                        StopMove();

                    if (target != null && IsInAttackRange(target))
                    {
                        ChangeRoleState(HallRoleState.Attack);
                        TryAttack(target);
                    }
                }
                else if(curRoleState == HallRoleState.Attack)
                {
                    if (target != null && target.IsAlive())
                    {
                        if (IsInAttackRange(target))
                            TryAttack(target);
                        else
                            ChangeRoleState(HallRoleState.Waiting);
                    }
                }
            }
        }

        protected bool CheckWin()
        {
            if (targets == null || targets.Count == 0)
                return false;

            for (int i = 0; i < targets.Count; i++)
            {
                if (targets[i] != null && targets[i].IsAlive())
                {
                    //Debug.LogError($"{targets[i].gameObject} is alive");
                    return false;
                }
            }

            ChangeRoleState(HallRoleState.WaitingSettlement);
            return true;
        }

        // 获取最近目标
        protected HallMonster GetNearestTarget()
        {
            if (targets == null || targets.Count == 0)
                return null;

            HallMonster target = null;
            var minDis = 9999f;
            for (int i = 0; i < targets.Count; i++)
            {
                if (!targets[i].IsAlive())
                    continue;

                var dis = Vector2.Distance(transform.localPosition, targets[i].transform.localPosition);
                if (dis < minDis)
                {
                    target = targets[i];
                    minDis = dis;
                }
            }

            return target;
        }

        // 检查是否在攻击范围内
        protected virtual bool IsInAttackRange(HallMonster target)
        {
            if (target == null) return false;

            //Debug.LogError($"Role: {transform.localPosition.x} {target.transform.localPosition.x}, {attackRange}");
            return (target.transform.localPosition.x - transform.localPosition.x <= attackRange) || 
                (Vector2.Distance(transform.localPosition, target.transform.localPosition) <= attackRange);
        }
        #endregion

        #region 状态切换

        public virtual bool IsAlive()
        {
            return !IsState(HallRoleState.Die);
        }

        public virtual bool IsMoving()
        {
            return curSpineName == GameConst.ModelAnim.run.ToString();
        }

        public virtual bool IsState(HallRoleState state)
        {
            return curRoleState == state;
        }

        public virtual void DoSpecialSkill(int cfgId)
        {
            if (curRoleState == HallRoleState.Attack)
            {
                specialSkillCfgID = cfgId;
                lastAttackTime = 0;
            }
        }

        public virtual void ChangeRoleState(HallRoleState state)
        {
            if (curRoleState == state)
                return;

            //Debug.LogError("Cur Role State is " + state);
            curRoleState = state;
            switch (state)
            {
                case HallRoleState.Rest:
                    Rest();
                    break;

                case HallRoleState.RestToBattle: // 休息室去战场的转场
                    RestToBattle();
                    break;

                case HallRoleState.MoveToBattlefield: // 跑向战场
                    MoveToBattlefield();
                    break;

                case HallRoleState.Waiting: // 等待怪物
                    Waiting();
                    break;

                case HallRoleState.Attack:// 攻击
                    Attack();
                    break;

                case HallRoleState.TakingDamage: // 受击
                    break;

                case HallRoleState.Revive: // 复活
                    Revive();
                    break;

                case HallRoleState.Die:
                    Die();
                    break;

                case HallRoleState.WaitingSettlement: // 等待结算
                    WaitingSettlement();
                    break;
                    
            }
        }

        // 休息室
        public virtual void Rest()
        {
            Clear();
            InitPosition(Vector2.zero);
            PlayLoopSpine(GameConst.ModelAnim.stand.ToString());
        }

        // 休息室跑向战场
        public virtual void RestToBattle()
        {
            DoMove(2);
            attackTimes = 0;
        }

        // 跑向战场
        public virtual void MoveToBattlefield()
        {
            InitData();
            EnableBoxCollider();
            DoMove(2f);
        }

        // 角色跑动到战场边缘，等待怪物
        public virtual void Waiting()
        {
            lastAttackTime = 0;
            if (transform.localPosition.x >= limitX)
                StopMove();
            PlayLoopSpine(GameConst.ModelAnim.run.ToString());

            PassiveMove(true);
        }

        // 攻击
        public virtual void Attack()
        {
            StopMove();
            PlayLoopSpine(GameConst.ModelAnim.stand.ToString());

            PassiveMove(false);
        }

        // 复活，下一关
        public virtual void Revive()
        {
            InitData();
            if (transform.localPosition.x < limitX)
                DoMove();

            PlayLoopSpine(GameConst.ModelAnim.run.ToString());

            EnableBoxCollider();
            ChangeRoleState(HallRoleState.Waiting);
        }

        // 死亡
        public virtual void Die()
        {
            Clear();
            curSpineName = GameConst.ModelAnim.dead.ToString();
            SpineUtil.Instance.PlayOnce(spineTrans, curSpineName, null, timeScale);
        }

        // 等待结算
        public virtual void WaitingSettlement()
        {
            //Debug.LogError("WaitingSettlement IsRoleWin = true ");
            UniGameEvent.Send(new GameEvent.HallBattleSettlement() { IsRoleWin = true });
            Clear();
            PlayLoopSpine(GameConst.ModelAnim.stand.ToString());

            PassiveMove(false);
        }

        private void Clear()
        {
            lastAttackTime = 0;
            targets.Clear();
            StopMove();
            DisableBoxCollider();
        }
        #endregion

        #region 攻击
        // 尝试攻击
        protected virtual void TryAttack(HallMonster target)
        {
            if (Time.time - lastAttackTime >= attackCooldown)
            {
                DoAttack(target);
            }
        }

        // 攻击逻辑
        public virtual void DoAttack(HallMonster target)
        {
            attackTimes++;
            lastAttackTime = Time.time;

            PassiveMove(false);
        }
        #endregion

        #region 碰撞检测
        private void OnTriggerEnter2D(Collider2D other)
        {
            // 只有攻击碰撞体激活时才检测碰撞
            var effectInfo = other.gameObject.GetComponent<HallEffectInfo>();
            if(effectInfo != null && effectInfo.hallEffectType == HallEffectType.Monster)
            {
                //Debug.LogError($"role be hit {effectInfo.AttackPower}, {currentHealth}");
                ShowBeHitHPChangeText("-", (int)effectInfo.AttackPower).Forget();
                currentHealth -= effectInfo.AttackPower;
                if (currentHealth <= 0)
                    UniGameEvent.Send(new GameEvent.HallBattleSettlement() { IsRoleWin = false });
                else
                    ShowRoleBeHitEffect(false);
            }
        }

        protected virtual async UniTask ShowBeHitHPChangeText(string str, int attack)
        {
            var poolKey = "Hall_HP_Change_Role_Text";

            if (!PoolHelper.Instance.CheckPool(poolKey))
                PoolHelper.Instance.AddPool(poolKey, hpChangeGo);

            await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));

            var pos = transform.localPosition;
            pos.y += 150;
            var go = PoolHelper.Instance.GetPoolObj(poolKey);
            go.transform.SetParent(hpChangeGo.transform.parent);
            go.transform.localPosition = pos;
            go.transform.localScale = Vector3.one;

            var endPos = pos;
            endPos.y += 120f;
            endPos.x -= 240;

            var midPos = (pos + endPos) * 0.5f;
            midPos.y += 50;

            var paths = new Vector3[] { pos, midPos, endPos };

            var text = go.GetComponent<Text>();
            text.text = str + attack;
            changeColor(100);
            go.SetActive(true);

            var actTime = 0.5f / timeScale;
            go.transform.DOLocalPath(paths, actTime, pathType: PathType.CatmullRom);
            DOTween.To(() => 100, changeColor, 30, actTime).SetEase(Ease.Linear).OnComplete(() =>
            {
                if (go == null)
                    return;

                go.SetActive(false);
                PoolHelper.Instance.RecycObj(go);
            });

            void changeColor(int value)
            {
                var color = text.color;
                color.a = value * 0.01f;
                text.color = color;
            }
        }

        private void ShowRoleBeHitEffect(bool bossMasterStroke)
        {
            var animName = bossMasterStroke ? "UIAttack_2" : "UI_Hurt";
            PoolHelper.Instance.LoadObj(animName, GameConst.GetEffPath(animName), (go) =>
            {
                if (curRoleState != HallRoleState.Attack)
                {
                    PoolHelper.Instance.RecycObj(go);
                    return;
                }

                go.transform.SetParent(transform);
                go.transform.localPosition = Vector3.zero;
                go.transform.localScale = Vector3.one * (bossMasterStroke ? 1.5f : 0.5f);
                go.SetActive(true);

                if (bossMasterStroke)
                    SpineUtil.Instance.PlayLoop(go.transform.Find("spine"), "hit", timeScale);

                GameHelper.Instance.SetDeayTimeDo(() =>
                {
                    if (go != null)
                        PoolHelper.Instance.RecycObj(go);
                }, 0.5f / timeScale);
            });
        }
        #endregion

        #region Common
        protected virtual void PassiveMove(bool value)
        {
            // 怪物、背景停止被动移动
            UniGameEvent.Send(new GameEvent.HallBattlePassiveMove() { passive = value });
        }

        public virtual void SuddenDeath()
        {
            ChangeRoleState(HallRoleState.Die);
        }

        protected virtual void DoMove(float add = 1)
        {
            if (rb != null)
                rb.velocity = moveDir * moveSpeed * timeScale * add;

            PlayLoopSpine(GameConst.ModelAnim.run.ToString());
        }

        protected virtual void StopMove()
        {
            if (rb != null)
                rb.velocity = Vector2.zero;
        }

        protected virtual void EnableBoxCollider()
        {
            if (boxCollider != null)
                boxCollider.enabled = true;
        }

        protected virtual void DisableBoxCollider()
        {
            if (boxCollider != null)
                boxCollider.enabled = false;
        }

        protected virtual void PlayLoopSpine(string name)
        {
            if (curSpineName == name)
                return;

            curSpineName = name;

            if (spineTrans == null) return;

            SpineUtil.Instance.PlayLoop(spineTrans, name, timeScale);
        }

        protected virtual void SetEffectSpineTimeScale(GameObject go)
        {
            var skeletonGraphic = go.GetComponentInChildren<SkeletonGraphic>();
            if (skeletonGraphic != null)
                skeletonGraphic.AnimationState.TimeScale = timeScale;

            var skeletonAnimation = go.GetComponentInChildren<SkeletonAnimation>();
            if (skeletonAnimation != null)
                skeletonAnimation.state.TimeScale = timeScale;
        }
        #endregion
    }
}