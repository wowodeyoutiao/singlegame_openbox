using Cysharp.Threading.Tasks;
using DG.Tweening;
using HotUpdateScripts;
using JEngine.Core;
using JEngine.Helper;
using Spine.Unity;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using static GameEvent;

namespace HallBattle
{
    public enum HallMonsterState
    {
        None = 0,
        MoveToBattlefield, // 跑向战场
        Pause, // 暂停
        Attack, // 攻击
        TakingDamage, // 受击
        Die, // 死亡
        WaitingSettlement, // 等待结算
    }

    public class HallBaseMonster : MonoBehaviour
    {
        private class HpChangeInfo
        {
            public int attackPower;
            public float hitDelayTime;
            public HallHPChangeType hPChangeType;
            public string hitEffectName;
            public Vector3 hitEffectPos;

            private float startTime;
            
            public HpChangeInfo(int power, float dTime, HallHPChangeType hPChangeType, string effectName, Vector3 pos)
            {
                attackPower = power;
                hitDelayTime = dTime;
                this.hPChangeType = hPChangeType;
                hitEffectName = effectName;
                hitEffectPos = pos;
            }

            public void StartWaiting()
            {
                startTime = Time.time;
            }

            public void ActImmediately()
            {
                startTime = Time.time;
                hitDelayTime = 0;
            }

            public bool CanShow()
            {
                return startTime + hitDelayTime <= Time.time;
            }
        }

        [Header("表")]
        protected int id;
        protected CFGItem_Monster cFGItem_Monster;

        [Header("基本属性")]
        protected bool isBoss; // 是否是boss关
        protected float scale;
        protected float maxHealth = 100f; // 初始血量
        protected float currentHealth; // 当前血量
        protected int attackPower = 10; // 战力
        protected int beHitTimes; // 受击次数
        protected int seckillTimes; // x次后直接秒杀英雄
        protected bool isMopUp; // 扫荡
        protected float limitX; // 战斗位置，最多可移动的地方

        protected RectTransform rectTransform;
        protected bool isMove;
        protected float moveSpeed = 5.2f; // 移速
        protected float attackRange = 300f; // 攻击范围

        protected float timeScale = 1f; // 攻速
        protected float attackCooldown = 2f;// 攻击冷却
        protected float lastAttackTime; // 上次攻击时间

        protected Vector2 moveDir = new Vector2(-1, 0);

        [Header("战斗状态")]
        protected HallMonsterState preMonsterState; // 角色状态
        protected HallMonsterState curMonsterState; // 角色状态

        [Header("使用技能")]
        protected bool useSkill;
        protected int skillIndex = 0;

        protected Rigidbody2D rb;
        protected BoxCollider2D boxCollider;
        protected bool IsRomate;
        protected string curSpineName;
        protected List<HallRole> targets = new List<HallRole>(); // 攻击目标

        // 血条
        protected GameObject hpHudGo;
        protected Slider hpSlider;

        // 飞血容器
        private GameObject[] hpChangeGos = new GameObject[2];

        // 假的移动
        protected bool passiveMove;

        // 死亡xs后，消失
        private float dieTimeInterval;

        // 多次掉血
        private List<HpChangeInfo> hpChangeInfos = new List<HpChangeInfo>();

        protected bool isPause;

        protected virtual void Awake()
        {
            LoadRB();
            limitX = 30;

            UniGameEvent.AddEventListener<GameEvent.HallBallPauseOrPlay>(HallBallPauseOrPlayEvt);
        }

        protected virtual void OnDestroy()
        {
            UniGameEvent.RemoveEventListener<GameEvent.HallBallPauseOrPlay>(HallBallPauseOrPlayEvt);
        }

        private void HallBallPauseOrPlayEvt(HallBallPauseOrPlay arge)
        {
            if (curMonsterState == HallMonsterState.Die)
                return;

            isPause = arge.pause;
            if(isPause)
            {
                preMonsterState = curMonsterState;
                ChangeMonsterState(HallMonsterState.Pause);
            }else if(preMonsterState != HallMonsterState.None)
            {
                ChangeMonsterState(preMonsterState);
            }
        }

        private void LoadRB()
        {
            if (rb == null)
            {
                rb = gameObject.AddComponent<Rigidbody2D>();
                rb.bodyType = RigidbodyType2D.Kinematic;

                boxCollider = gameObject.AddComponent<BoxCollider2D>();
                boxCollider.isTrigger = true;
                boxCollider.offset = new Vector2(100, 150);
                boxCollider.size = new Vector2(200, 400);
            }

            if (rectTransform == null)
                rectTransform = gameObject.GetComponent<RectTransform>();
        }

        // 初始化（包含复用的）
        public virtual void InitWithData(int[] info, Vector2 birthPosition, bool isBoss, float scale, bool isMopUp)
        {
            id = info[0];
            maxHealth = info[1];
            currentHealth = info[1];
            attackPower = info[2];
            seckillTimes = info[3];
            this.isBoss = isBoss;
            this.scale = scale;
            this.isMopUp = isMopUp;

            cFGItem_Monster = CFG_Monster.GetById(id);
            useSkill = (cFGItem_Monster.Skills != null && cFGItem_Monster.Skills.Length > 1);
            skillIndex = 0;
            IsRomate = cFGItem_Monster.Model.ToLower().Contains("romate");

            if (!IsRomate && cFGItem_Monster.AttackDistance <= 3)
                attackRange = 220 + Random.Range(-20, 50);
            else
                attackRange = 360 + Random.Range(-20, 50);

            attackCooldown = Random.Range(1.8f, 2.5f);

            transform.localPosition = birthPosition;
            transform.localScale = Vector3.one * scale;

            LoadRB();
            if (rb != null)
                rb.velocity = Vector2.zero;

            isMove = false;
            preMonsterState = HallMonsterState.None;
            isPause = false;
            passiveMove = false;
            beHitTimes = 0;
            dieTimeInterval = 0;
            curSpineName = string.Empty;
            targets.Clear();
            hpChangeInfos.Clear();
        }

        public virtual async UniTask AddHPHud(GameObject go)
        {
            if(isBoss)
            {
                if (hpHudGo != null)
                    hpHudGo.SetActive(false);
                return;
            }

            if (hpHudGo == null)
            {
                var poolKey = "Hall_Monster_HP_Hud";
                if (!PoolHelper.Instance.CheckPool(poolKey))
                    PoolHelper.Instance.AddPool(poolKey, go);

                await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
                hpHudGo = PoolHelper.Instance.GetPoolObj(poolKey);
                hpSlider = hpHudGo.GetComponentInChildren<Slider>();
            }

            hpHudGo.transform.SetParent(transform);
            hpHudGo.transform.localScale = Vector3.one / scale;
            hpHudGo.transform.localPosition = new Vector3(0, 446, 0);
            hpHudGo.SetActive(true);

            hpSlider.value = 1;
        }

        public virtual void AddHpChangeTrans(Transform trans)
        {
            hpChangeGos[0] = trans.GetChild(0).gameObject;
            hpChangeGos[1] = trans.GetChild(1).gameObject;
        }

        #region 战场战斗
        public virtual void AddTarget(HallRole hallRole)
        {
            targets.Add(hallRole);
        }

        protected virtual void Update()
        {
            if (isPause) return;

            if(curMonsterState == HallMonsterState.MoveToBattlefield || curMonsterState == HallMonsterState.Attack)
            {
                var target = GetNearestTarget();
                if (curMonsterState == HallMonsterState.MoveToBattlefield)
                {
                    if (target != null && IsInAttackRange(target))
                        ChangeMonsterState(HallMonsterState.Attack);
                }
                else if (curMonsterState == HallMonsterState.Attack)
                {
                    if (target != null && IsInAttackRange(target))
                        TryAttack(target);
                }
            }

            if (passiveMove)
                DoMove();

            if(curMonsterState == HallMonsterState.Die)
            {
                if(dieTimeInterval > 1.5f)
                    gameObject.SetActive(false);

                dieTimeInterval += Time.deltaTime;
            }

            CheckShowBeHitHPChange();
        }

        protected virtual void FixedUpdate()
        {
            if (isPause) return;
            if (isMove)
            {
                if (targets != null && targets.Count > 0 && targets[0].transform != null &&
                    transform.localPosition.x - targets[0].transform.localPosition.x <= 160 && curMonsterState != HallMonsterState.Die)
                {
                    StopMove();
                    PlayLoopSpine(GameConst.ModelAnim.stand.ToString());
                    return;
                }

                var pos = transform.localPosition;
                pos.x -= moveSpeed * timeScale;
                rectTransform.anchoredPosition = pos;
            }
        }

        // 获取最近目标
        protected HallRole GetNearestTarget()
        {
            if (targets == null || targets.Count == 0)
                return null;

            HallRole target = null;
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
        #endregion

        #region 状态切换
        public virtual void ChangeMonsterState(HallMonsterState state)
        {
            if (curMonsterState == state)
                return;

            //Debug.LogError($"ChangeMonsterState: {gameObject.name}, {curMonsterState}, {state}");
            curMonsterState = state;
            switch (state)
            {
                case HallMonsterState.MoveToBattlefield: // 跑向战场
                    MoveToBattlefield();
                    break;

                case HallMonsterState.Pause:
                    Pause();
                    break;

                case HallMonsterState.Attack:// 攻击
                    Attack();
                    break;

                case HallMonsterState.TakingDamage: // 受击
                    TakingDamage();
                    break;

                case HallMonsterState.Die:
                    Die();
                    break;

                case HallMonsterState.WaitingSettlement:
                    WaitingSettlement();
                    break;
            }
        }

        // 跑向战场
        public virtual void MoveToBattlefield()
        {
            EnableBoxCollider();
            DoMove();
            PlayLoopSpine(GameConst.ModelAnim.walk.ToString());
        }

        public virtual void Pause()
        {
            PlayLoopSpine(GameConst.ModelAnim.stand.ToString());
        }

        // 攻击
        public virtual void Attack()
        {
            var target = GetNearestTarget();
            //Debug.LogError($"Attack  {passiveMove}, {target.IsMoving()}");
            if (target != null && target.IsMoving())
            {
                DoMove();
                PlayLoopSpine(GameConst.ModelAnim.walk.ToString());
            }
            else
            {
                StopMove();
                PlayLoopSpine(GameConst.ModelAnim.stand.ToString());
            }
        }

        public virtual void TakingDamage()
        {
            PlayLoopSpine(GameConst.ModelAnim.stand.ToString());
        }

        // 死亡
        public virtual void Die()
        {
            Clear();
            curSpineName = GameConst.ModelAnim.dead.ToString();
            SpineUtil.Instance.PlayOnce(transform, curSpineName, null, timeScale);
        }

        public virtual void WaitingSettlement()
        {
            DisableBoxCollider();
        }

        private void Clear()
        {
            if (hpHudGo != null)
                hpHudGo.SetActive(false);
            targets.Clear();
            StopMove();
            DisableBoxCollider();
        }
        #endregion

        #region 攻击
        // 检查是否在攻击范围内
        protected virtual bool IsInAttackRange(HallRole target)
        {
            if (target == null) return false;
            //Debug.LogError($"Monster: {transform.localPosition.x} {target.transform.localPosition.x}, {attackRange}");
            return Vector2.Distance(transform.localPosition, target.transform.localPosition) <= attackRange;
        }

        // 尝试攻击
        protected virtual void TryAttack(HallRole target)
        {
            if (Time.time - lastAttackTime >= attackCooldown / timeScale)
                Attack(target);
        }

        // 攻击逻辑
        public virtual void Attack(HallRole target)
        {
            lastAttackTime = Time.time;

            var animName = "attack";
            if (useSkill && skillIndex != 0)
            {
                animName = "skill";
                skillIndex = 0;
            }
            else
                skillIndex++;
            curSpineName = animName;
            SpineUtil.Instance.PlayOnce(transform, animName, () =>
            {
                PlayLoopSpine(GameConst.ModelAnim.stand.ToString());
            }, timeScale);
        }
        #endregion

        #region 受伤检测
        private void OnTriggerEnter2D(Collider2D other)
        {
            if (curMonsterState == HallMonsterState.Die)
                return;

            // 只有攻击碰撞体激活时才检测碰撞
            var effectInfo = other.gameObject.GetComponent<HallEffectInfo>();
            if (effectInfo != null && effectInfo.hallEffectType == HallEffectType.Role)
            {
                if (effectInfo.AttackNum <= 0 && !effectInfo.CanHit(gameObject.name))
                    return;

                effectInfo.ChangeAttackNum(-1);
                effectInfo.AddWounded(gameObject.name);
                //Debug.LogError($"monster be hit:{gameObject.name}, {curMonsterState}, {currentHealth}, {effectInfo.AttackPower}, {isBoss}");

                if (isBoss)
                    UniGameEvent.Send(new GameEvent.HallBattleBossBeHit() { attackPower = effectInfo.AttackPower });

                beHitTimes++;
                BeHit(effectInfo);
                SpecialSkillBuff(effectInfo.skillMainCfgID, other.transform.position);
            }
        }

        private void BeHit(HallEffectInfo effectInfo)
        {
            if (effectInfo.hitNum <= 1)
            {
                if (hpChangeInfos == null || hpChangeInfos.Count == 0)
                    ShowBeHitHPChange((int)effectInfo.AttackPower, effectInfo.hallHPChangeType, effectInfo.HitEffectName, effectInfo.HitEffectPos);
                else
                    hpChangeInfos.Add(new HpChangeInfo((int)effectInfo.AttackPower, 0, effectInfo.hallHPChangeType, effectInfo.HitEffectName, effectInfo.HitEffectPos));
            }
            else
            {
                var actImmediately = hpChangeInfos == null || hpChangeInfos.Count == 0;
                var singlePower = effectInfo.AttackPower / effectInfo.hitNum;
                var singleTime = effectInfo.hitTime / effectInfo.hitNum;
                var totalPower = 0f;
                var totalTime = 0f;
                for(int i = 0; i < effectInfo.hitNum; i++)
                {
                    var power = 0f;
                    var dTime = 0f;
                    if(i == effectInfo.hitNum - 1)
                    {
                        power = effectInfo.AttackPower - totalPower;
                        dTime = effectInfo.hitTime - totalTime;
                    }
                    else
                    {
                        power = (Random.Range(-10, 10) * 0.01f + 1) * singlePower;
                        totalPower += power;

                        dTime = (Random.Range(-10, 10) * 0.01f + 1) * singleTime;
                        totalTime += dTime;
                    }

                    hpChangeInfos.Add(new HpChangeInfo((int)power, dTime, effectInfo.hallHPChangeType, effectInfo.HitEffectName, effectInfo.HitEffectPos));
                }

                if(actImmediately)
                {
                    hpChangeInfos[0].ActImmediately();
                }
            }
        }

        private void ShowBeHitHPChange(int attack, HallHPChangeType hPChangeType, string hitEffectName, Vector3 pos)
        {
            ShowBeHitHPChangeText(attack, hPChangeType).Forget();

            currentHealth -= attack;

            if (!isBoss && hpSlider != null)
                hpSlider.value = currentHealth / maxHealth;

            if (currentHealth <= 0)
            {
                // 发消息，掉落武器
                UniGameEvent.Send(new GameEvent.HallBattleDropItem() { position = transform.localPosition });
                ChangeMonsterState(HallMonsterState.Die);
            }

            if (!string.IsNullOrEmpty(hitEffectName))
                ShowBeHitEffect(hitEffectName, pos).Forget();
        }

        private void CheckShowBeHitHPChange()
        {
            if (curMonsterState == HallMonsterState.Die)
                return;

            if (hpChangeInfos == null || hpChangeInfos.Count == 0)
                return;

            if (hpChangeInfos[0].CanShow())
            {
                ShowBeHitHPChange((int)hpChangeInfos[0].attackPower, hpChangeInfos[0].hPChangeType, hpChangeInfos[0].hitEffectName, hpChangeInfos[0].hitEffectPos);
                hpChangeInfos.RemoveAt(0);

                if (hpChangeInfos != null && hpChangeInfos.Count > 0)
                    hpChangeInfos[0].StartWaiting();
            }
        }

        protected virtual async UniTask ShowBeHitHPChangeText(int attack, HallHPChangeType hPChangeType)
        {
            var poolKey = "Hall_HP_Change_Text" + hPChangeType;
            GameObject go;
            if (hPChangeType == HallHPChangeType.FS)
                go = hpChangeGos[0];
            else
                go = hpChangeGos[1];

            if (!PoolHelper.Instance.CheckPool(poolKey))
                PoolHelper.Instance.AddPool(poolKey, go);

            await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));

            var pos = transform.localPosition;
            pos.y += 150;
            var hpChangeGo = PoolHelper.Instance.GetPoolObj(poolKey);
            hpChangeGo.transform.SetParent(go.transform.parent);
            hpChangeGo.transform.localPosition = pos;
            hpChangeGo.transform.localScale = Vector3.one * 0.1f;

            var text = hpChangeGo.GetComponent<Text>();
            text.text = "-" + attack;
            changeColor(100);
            hpChangeGo.SetActive(true);

            var targetPos = pos;
            targetPos.y += 260;
            targetPos.x += Random.Range(-360, 360);
            var timeInterval = 0.7f;
            hpChangeGo.transform.DOLocalMove(targetPos, timeInterval).SetEase(Ease.Linear);
            DOTween.To(() => 100, changeColor, 30, timeInterval).SetEase(Ease.Linear).OnComplete(() =>
            {
                if (hpChangeGo == null)
                    return;

                hpChangeGo.SetActive(false);
                PoolHelper.Instance.RecycObj(hpChangeGo);
            });

            void changeColor(int value)
            {
                var color = text.color;
                color.a = value * 0.01f;
                text.color = color;
            }
        }

        private async UniTask ShowBeHitEffect(string effectName, Vector3 pos)
        {
            // 技能
            if (!PoolHelper.Instance.CheckPool(effectName))
            {
                var baseObj = new JPrefab(GameConst.GetEffPath(effectName), true);
                await UniTask.WaitUntil(() => baseObj.Loaded);

                PoolHelper.Instance.AddPool(effectName, baseObj.Instance);
            }
            await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(effectName));
            var effGo = PoolHelper.Instance.GetPoolObj(effectName);
            effGo.transform.SetParent(transform);
            effGo.transform.localPosition = pos;
            effGo.transform.localScale = Vector3.one;
            effGo.SetActive(true);
            SetEffectSpineTimeScale(effGo);

            GameHelper.Instance.SetDeayTimeDo(() =>
            {
                if (effGo == null || PoolHelper.Instance == null)
                    return;
                PoolHelper.Instance.RecycObj(effGo);
            }, 1f / timeScale);
        }

        private void SpecialSkillBuff(int specialSkillID, Vector3 centerPos)
        {
            if (specialSkillID == 0 || curMonsterState == HallMonsterState.Die)
                return;

            if (specialSkillID == 12)
            {
                // 碎星地裂斩,将怪物退走
                preMonsterState = curMonsterState;
                ChangeMonsterState(HallMonsterState.TakingDamage);

                transform.DOLocalMoveX(350, 1.2f).OnComplete(() =>
                {
                    if (curMonsterState == HallMonsterState.TakingDamage)
                        ChangeMonsterState(HallMonsterState.MoveToBattlefield);
                });
            }
            else if (specialSkillID == 13)
            {
                // 虚空坍缩,将周围怪物吸入并造成伤害
                var cfg = CFG_SkillMainBattle.GetById(specialSkillID);
                if (cfg == null || cfg.HitPositionInfo == null || cfg.HitPositionInfo.Length < 2)
                    return;

                var dir = (centerPos - transform.position).normalized;
                DOTween.To(() => 0, (a) => 
                {
                    var dis = Vector2.Distance(transform.position, centerPos);
                    if(dis > 0.1f)
                    {
                        var pos = transform.position;
                        pos += dir * 0.01f;
                        transform.position = pos;
                    }
                
                }, 80, 0.2f).SetEase(Ease.Linear);
            }
            else if (specialSkillID == 16)
            {
                // 荆棘牢笼,受到伤害，且不可移动、攻击
                preMonsterState = curMonsterState;
                ChangeMonsterState(HallMonsterState.TakingDamage);

                GameHelper.Instance.SetDeayTimeDo(() =>
                {
                    if (curMonsterState == HallMonsterState.TakingDamage)
                        ChangeMonsterState(preMonsterState);
                }, 2f);
            }
            else if (specialSkillID == 18)
            {
                // 时空之轮，伤害、减速
                timeScale = 0.5f;
                GameHelper.Instance.SetDeayTimeDo(() =>
                {
                    timeScale = 1;
                }, 2f);
            }
        }
        #endregion

        #region Common
        public virtual bool IsAlive()
        {
            return curMonsterState != HallMonsterState.Die;
        }

        public virtual void PassiveMove(bool value)
        {
            passiveMove = value;
            if (!value && curMonsterState != HallMonsterState.MoveToBattlefield)
                StopMove();
        }

        protected virtual void DoMove()
        {
            moveSpeed = (passiveMove && curMonsterState != HallMonsterState.MoveToBattlefield) ? 5 : 5.2f;
            isMove = true;
        }

        protected virtual void StopMove()
        {
            isMove = false;
        }

        protected virtual void EnableBoxCollider()
        {
            if (boxCollider != null)
                boxCollider.enabled = true;
        }

        protected virtual void DisableBoxCollider()
        {
            //Debug.LogError("DisableBoxCollider");
            if (boxCollider != null)
                boxCollider.enabled = false;
        }

        protected virtual void PlayLoopSpine(string name)
        {
            if (curSpineName == name)
                return;

            curSpineName = name;
            SpineUtil.Instance.PlayLoop(transform, name, timeScale);
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