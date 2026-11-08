using Cysharp.Threading.Tasks;

using DG.Tweening;

using HotUpdateScripts.ProtoM;

using JEngine.Core;
using JEngine.Helper;

using pb;

using Spine.Unity;

using System;
using System.Collections.Generic;

using UnityEngine;

namespace HotUpdateScripts.Game
{
    public class TurnBasedRole : MonoBehaviour
    {
        private Transform m_Spine;

        private Ease m_MoveEase = Ease.OutSine;

        private CFGItem_Skill m_SkillCfg;

        private CFGItem_SkillEffect m_SkillEff;

        // 攻击时是否需要移动到对手面前
        private bool m_IsNeedMove = true;

        // 已经离开自己的位置了？
        private bool m_IsMoved;

        private const float m_BuffTime = 0.7f;

        private const float m_EndWaitTime = 0.2f;
        // 怪物远程攻击特效弹道时间
        private const float m_MonsterEffTime = 0.08f;

        private Dictionary<string, GameObject> m_BuffObjList = new Dictionary<string, GameObject>();

        public SkeletonAnimation spineAnim;
        public GameConst.ModelAnim animState;

        public ObjectType roleType;
        public TurnBasedBattle.ActionCamp campType;

        //public UnityAction OnComplete;

        public Transform root_hub;

        private Tween m_Tween;

        private List<int> m_BufferList = new List<int>();
        public List<int> bufferList => m_BufferList;

        public long uid;

        public int cfgId;
        /// <summary>
        /// 站位 （1 - 6）
        /// </summary>
        public int position;

        /// <summary>
        /// 角色名字
        /// </summary>
        public string roleName;

        /// <summary>
        /// 最大血量
        /// </summary>
        public int maxHp;

        /// <summary>
        /// 当前血量
        /// </summary>
        public int currentHp;

        /// <summary>
        /// 技能ID
        /// </summary>
        public int skillId;

        /// <summary>
        /// 已死亡
        /// </summary>
        public bool isDead;

        /// <summary>
        /// 职业
        /// </summary>
        public int job;

        /// <summary>
        /// 主元素。
        /// 怪物走配置。
        /// </summary>
        public int mainElement;

        /// <summary>
        /// 目标位子类型。
        /// Movement 近战攻击时需要移动到角色面前
        /// Effect 特效播放位子，一般在角色身上
        /// Buff 特效一般在角色身上
        /// </summary>
        public enum GetPositionType
        {
            Movement,
            Effect,
            Buff,
        }


        private void Awake()
        {
            m_Spine = transform.Find("spine");
            if (m_Spine)
                spineAnim = m_Spine.GetComponent<SkeletonAnimation>();

            root_hub = transform.Find("root_hub");
        }

        private void Start()
        {
            Init();
        }

        public void Init()
        {
            animState = GameConst.ModelAnim.No;
            ChangeAnimationState(GameConst.ModelAnim.stand);
        }

        public void ChangeAnimationState(GameConst.ModelAnim actName, Action callback = null)
        {
            if (actName != GameConst.ModelAnim.dead && actName != GameConst.ModelAnim.attack && !animState.ToString().Contains("skill") && animState.CompareTo(actName) == 0)
                return;
            animState = actName;

            if (spineAnim == null || (roleType != ObjectType.PLAYER && roleType != ObjectType.PARTNER && roleType != ObjectType.MONSTER))
                return;

            if (!IsContainsAnimation(animState.ToString()))
            {
                Debug.Log($"<color=red>spine 动画名 {animState} 不存在!  .. {roleName}<>");
                return;
            }

            if (animState == GameConst.ModelAnim.attack || animState.ToString().Contains("skill"))
                SpineUtil.Instance.PlayOnce(spineAnim, animState.ToString(), () =>
                {
                    animState = GameConst.ModelAnim.stand;
                    SpineUtil.Instance.PlayLoop(spineAnim, animState.ToString());
                }, 3);
            else if (actName == GameConst.ModelAnim.dead)
                SpineUtil.Instance.PlayOnce(spineAnim, animState.ToString(), callback, 3);
            else
                SpineUtil.Instance.PlayLoop(spineAnim, animState.ToString());
        }

        public void StartAttack()
        {
            Debug.Log($"**** attack uid：{uid}， name: = {roleName}，skillId = {skillId}, attackNum = {TurnBasedFightCtrl.Instance.attackNum}");
            TurnBasedFightCtrl.Instance.AddBuffers();

            Sequence seq = DOTween.Sequence();
            if (skillId != 0)
            {
                m_SkillCfg = CFG_Skill.GetById(skillId);

                if (roleType == ObjectType.PLAYER || roleType == ObjectType.PARTNER)
                {
                    m_SkillEff = CFG_SkillEffect.GetById(m_SkillCfg.SkillID);
                    m_IsNeedMove = IsNeedMove();
                    if (m_IsNeedMove)
                    {
                        seq.Append(MoveToTargetPosition());
                    }
                    seq.AppendCallback(RoleAttack);
                }
                else if (roleType == ObjectType.MONSTER)
                {
                    m_IsNeedMove = !CFG_Monster.IsRomate(cfgId);
                    if (m_IsNeedMove)
                    {
                        seq.Append(MoveToTargetPosition());
                    }
                    seq.AppendCallback(MonsterAttack);
                }
            }
            else
            {
                CoroutineCall(m_EndWaitTime, () =>
                {
                    TurnBasedFightCtrl.Instance.ChangeState(TurnBasedFightCtrl.BattleState.AttackCompleted);
                });
            }

            m_Tween = seq;
        }

        /// <summary>
        /// 被攻击目标的位置
        /// </summary>
        /// <param name="isOffset">角色移动时需要距离目标有段距离</param>
        /// <returns></returns>
        private Vector3 GetHitPositionOfStanding(GetPositionType positionType)
        {
            long uid = GetFirstLivingUid();

            if (campType == TurnBasedBattle.ActionCamp.BLUE_CAMP)
            {
                if (TurnBasedRoleManager.Instance.GetRedRoleByUid(uid, out TurnBasedRole role))
                {
                    return positionType switch
                    {
                        GetPositionType.Movement => role.GetWorldPosition() - new Vector3(2.3f, 0, 0),
                        GetPositionType.Effect => role.GetWorldPosition(),
                        GetPositionType.Buff => role.GetWorldPosition(),
                        _ => Vector3.zero,
                    };
                }
                else
                {
                    Debug.Log("未找到角色！ skill target role red uid " + uid);
                    return Vector3.zero;
                }
            }
            else if (campType == TurnBasedBattle.ActionCamp.RED_CAMP)
            {
                if (TurnBasedRoleManager.Instance.GetBlueRoleByUid(uid, out TurnBasedRole role))
                {
                    return positionType switch
                    {
                        GetPositionType.Movement => role.GetWorldPosition() + new Vector3(2.3f, 0, 0),
                        GetPositionType.Effect => role.GetWorldPosition(),
                        GetPositionType.Buff => role.GetWorldPosition(),
                        _ => Vector3.zero,
                    };
                }
                else
                {
                    Debug.Log("未找到角色！ skill target role blueuid " + uid);
                    return Vector3.zero;
                }
            }
            else
                return Vector3.zero;

        }

        public bool IsContainsAnimation(string animationName)
        {
            if (spineAnim == null)
                return false;

            foreach (var animation in spineAnim.skeleton.Data.Animations)
                if (animationName == animation.Name) return true;

            return false;
        }

        public Vector3 GetWorldPosition()
        {
            return transform.position;
        }

        private Tweener MoveToTargetPosition()
        {
            m_IsMoved = true;
            Vector3 targetPos = GetHitPositionOfStanding(GetPositionType.Movement);
            return transform.DOMove(targetPos, TurnBasedFightCtrl.Instance.attackMoveSpeed).SetEase(m_MoveEase);
        }

        // 移动到原位
        private Tweener MoveBackToStandPosition(Action onComplete = null)
        {
            return transform.DOLocalMove(Vector3.zero, TurnBasedFightCtrl.Instance.attackMoveSpeed).SetEase(m_MoveEase).OnComplete(() =>
            {
                m_IsMoved = false;
                onComplete?.Invoke();
            });
        }

        private bool IsNeedMove()
        {
            // 战士职业
            bool isWarrior = m_SkillCfg.SkillProfession % 10 == 1 && m_SkillCfg.Needaim == 1;
            // 直线技能特效（激光脉冲、地震波），需要移动到地方面前播放
            bool isLineSkillEffect = (m_SkillCfg.UUID >= 1701 && m_SkillCfg.UUID <= 1730) || (m_SkillCfg.UUID >= 2901 && m_SkillCfg.UUID <= 2930);
            return isWarrior || isLineSkillEffect;
        }

        // 玩家或伙伴攻击
        private void RoleAttack()
        {
            GameConst.ModelAnim animName = GameConst.GetAnim((int)m_SkillCfg.Needskill);
            Debug.Log("role skill aniName = " + animName);

            ChangeAnimationState(animName);
            PlayEffectRoleAtk().Forget();

            float waitTime = 0;
            if (!string.IsNullOrEmpty(m_SkillEff.StartEffect))
                waitTime = (float)m_SkillEff.StartRemainTime / (1000f * 4);
            else if (!string.IsNullOrEmpty(m_SkillEff.BallisticEffect))
                waitTime = (float)m_SkillEff.BallisticRemainTime / (1000f * 3);
            else
                waitTime = 0.1f;

            /*
            float bloodTextWaitTime = waitTime * 0.7f;
            if (m_SkillEff.TargetID == 16)
                bloodTextWaitTime = 0.5f;
            else if (m_SkillEff.TargetID == 17)
                bloodTextWaitTime = 0.2f;
            else if (m_SkillEff.TargetID == 29)
                bloodTextWaitTime = waitTime * 0.4f;*/

            CoroutineCall(waitTime, () =>
            {
                TurnBasedFightCtrl.Instance.UpdateRoleInfo();
                TurnBasedFightCtrl.Instance.CollectKilledRole();

                if (TurnBasedFightCtrl.Instance.IsNeedPlayHurtEffect(skillId))
                {
                    TurnBasedFightCtrl.Instance.PlayHurt(GetFirstLivingUid(), skillId);
                }
            });

            /*
            float hurtPlayWait = waitTime;
            if (m_SkillCfg.SkillProfession % 10 == 1)
                hurtPlayWait = waitTime * 0.5f;
            else if (m_SkillCfg.SkillProfession % 10 == 2)
                hurtPlayWait = waitTime * 0.9f;
            else if (m_SkillCfg.SkillProfession % 10 == 3)
                hurtPlayWait = waitTime * 0.9f;
            */
            //CoroutineCall(waitTime, () =>
            //{
            //    if (TurnBasedFightCtrl.Instance.IsNeedPlayHurtEffect(skillId))
            //    {
            //        TurnBasedFightCtrl.Instance.PlayHurt(GetFirstLivingUid(), skillId);
            //    }
            //});

            CoroutineCall(waitTime + 0.1f, () =>
            {
                float waitRevive = 0;
                if (TurnBasedFightCtrl.Instance.HasRevive())
                    waitRevive = 0.2f;

                RemoveUnusedBuffers();
                if (m_IsNeedMove)
                {
                    // 遇到出血技能死亡的时，角色已提前拉回原位，不需要再触发。
                    if (!isDead)
                        MoveBackToStandPosition();

                    CoroutineCall(TurnBasedFightCtrl.Instance.attackMoveSpeed + m_EndWaitTime + waitRevive, () =>
                    {
                        TurnBasedFightCtrl.Instance.ChangeState(TurnBasedFightCtrl.BattleState.AttackCompleted);
                    });
                }
                else
                {
                    CoroutineCall(m_EndWaitTime + waitRevive, () =>
                    {
                        TurnBasedFightCtrl.Instance.ChangeState(TurnBasedFightCtrl.BattleState.AttackCompleted);
                    });
                }
            });
        }

        /// <summary>
        /// 主角、伙伴播放技能特效
        /// </summary>
        public async UniTask PlayEffectRoleAtk()
        {
            if (CFG_Actionset.ChealkId(m_SkillCfg.SkillID) && !string.IsNullOrEmpty(CFG_Actionset.GetById(m_SkillCfg.SkillID).Music))
            {
                AudioManager.Instance.PlaySound(CFG_Actionset.GetById(m_SkillCfg.SkillID).Music);
            }

            if (!CFG_SkillEffect.Get().DicDatas.ContainsKey(m_SkillCfg.SkillID))
                return;

            Vector3 tpos = GetHitPositionOfStanding(GetPositionType.Effect);

            tpos.z = SceneManager.EffOffZ;

            if (!string.IsNullOrEmpty(m_SkillEff.BallisticEffect))
            {
                string poolKey = m_SkillEff.BallisticEffect;

                if (!PoolHelper.Instance.CheckPool(poolKey))
                {
                    var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
                    await UniTask.WaitUntil(() => baseObj.Loaded);
                    PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
                }
                await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
                var Eff = PoolHelper.Instance.GetPoolObj(poolKey);
                Eff.transform.SetParent(TurnBasedFightCtrl.Instance.EffectRoot);
                Eff.transform.position = new Vector3(transform.position.x, transform.position.y, SceneManager.EffOffZ);

                Vector2 direction = (transform.position - tpos);
                direction.Normalize();
                // 设置物体的朝向
                Eff.transform.rotation = Quaternion.AngleAxis(Mathf.Atan2(direction.y, direction.x) * Mathf.Rad2Deg, Vector3.forward);

                //CoroutineCall((float)m_SkillEff.BallisticDelay / (1000f * 3), () =>
                //{
                if (Eff == null)
                    return;
                Eff.SetActive(true);
                Eff.transform.DOMove(tpos, (float)m_SkillEff.BallisticRemainTime / (1000f * 2)).SetEase(Ease.Linear).OnComplete(() =>
                {
                    PoolHelper.Instance.RecycObj(Eff);
                });
                //});
            }

            if (!string.IsNullOrEmpty(m_SkillEff.StartEffect))
            {
                string poolKey = m_SkillEff.StartEffect;
                if (!PoolHelper.Instance.CheckPool(poolKey))
                {
                    var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
                    await UniTask.WaitUntil(() => baseObj.Loaded);

                    PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
                }
                await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
                var Eff = PoolHelper.Instance.GetPoolObj(poolKey);
                Eff.transform.SetParent(TurnBasedFightCtrl.Instance.EffectRoot);

                CoroutineCall((float)m_SkillEff.StartDelay / (1000f * 2), () =>
                {
                    if (Eff == null)
                        return;
                    Eff.transform.position = new Vector3(transform.position.x, transform.position.y, SceneManager.EffOffZ);

                    if (campType == TurnBasedBattle.ActionCamp.BLUE_CAMP)
                        Eff.transform.eulerAngles = new Vector3(0, 180, 0);
                    else if (campType == TurnBasedBattle.ActionCamp.RED_CAMP)
                        Eff.transform.eulerAngles = Vector3.zero;

                    Eff.SetActive(true);
                });

                CoroutineCall((float)m_SkillEff.StartRemainTime / 1000f, () =>
                    {
                        if (Eff == null)
                            return;
                        PoolHelper.Instance.RecycObj(Eff);
                    });
            }

            if (!string.IsNullOrEmpty(m_SkillEff.BodyEffect))
            {
                string poolKey = m_SkillEff.BodyEffect;
                if (!PoolHelper.Instance.CheckPool(poolKey))
                {
                    var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
                    await UniTask.WaitUntil(() => baseObj.Loaded);
                    PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
                }
                await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
                var Eff = PoolHelper.Instance.GetPoolObj(poolKey);
                Eff.transform.SetParent(TurnBasedFightCtrl.Instance.EffectRoot);
                Eff.transform.position = tpos;

                CoroutineCall((float)m_SkillEff.BodyDelay / (1000f * 2), () =>
                {
                    if (Eff == null)
                        return;
                    Eff.SetActive(true);
                    CoroutineCall((float)m_SkillEff.BodyRemainTime / 1000f, () =>
                    {
                        if (Eff == null)
                            return;
                        PoolHelper.Instance.RecycObj(Eff);
                    });
                });
            }
        }

        // 怪物攻击
        private void MonsterAttack()
        {
            ChangeAnimationState(GameConst.GetAnim((int)CFG_Skill.GetById(skillId).Needskill));

            GameConst.ModelAnim animName = GameConst.GetAnim((int)m_SkillCfg.Needskill);
            Debug.Log("monster skill aniName = " + animName);

            float waitTime = 0.5f;
            bool isRomate = CFG_Monster.IsRomate(cfgId);

            float updateHpTime, completeTime;

            if (isRomate)
            {
                CoroutineCall(waitTime * 0.5f, () =>
                {
                    PlayEffectMonsterAtk().Forget();
                });

                CoroutineCall(waitTime * 0.5f + m_MonsterEffTime, () =>
                {
                    if (TurnBasedFightCtrl.Instance.IsNeedPlayHurtEffect(skillId))
                    {
                        TurnBasedFightCtrl.Instance.PlayHurt(GetFirstLivingUid(), skillId);
                    }
                });

                updateHpTime = waitTime * 0.5f + m_MonsterEffTime;
                completeTime = waitTime + 0.1f;
            }
            else
            {
                updateHpTime = waitTime * 0.5f + TurnBasedFightCtrl.Instance.attackMoveSpeed;
                completeTime = waitTime + TurnBasedFightCtrl.Instance.attackMoveSpeed + 0.1f;

                CoroutineCall(waitTime, () =>
                {
                    MoveBackToStandPosition();
                });
            }

            CoroutineCall(updateHpTime, () =>
            {
                TurnBasedFightCtrl.Instance.UpdateRoleInfo();
                TurnBasedFightCtrl.Instance.CollectKilledRole();
            });

            float waitRevive = 0;
            if (TurnBasedFightCtrl.Instance.HasRevive())
                waitRevive = 0.2f;

            CoroutineCall(completeTime + waitRevive, () =>
            {
                TurnBasedFightCtrl.Instance.ChangeState(TurnBasedFightCtrl.BattleState.AttackCompleted);
            });
        }

        public async UniTask PlayEffectMonsterAtk()
        {
            var cfgInfo = CFG_Monster.GetById(cfgId);
            string poolKey = cfgInfo.AttackEffect;
            var SoundKey = AudioKey.Sound.BossAttack;
            string bulletIcon = "";

            if (cfgInfo.Isboss == 0 && !string.IsNullOrEmpty(poolKey))
            {
                if (poolKey.Contains("yg_"))
                {
                    poolKey = "Monster_zd";
                    bulletIcon = cfgInfo.AttackEffect;

                    SoundKey = string.IsNullOrEmpty(poolKey) ? AudioKey.Sound.JZAttack : AudioKey.Sound.YCAttack;
                }
            }

            AudioManager.Instance.PlaySound(SoundKey);

            if (string.IsNullOrEmpty(poolKey))
                return;

            var tpos = GetHitPositionOfStanding(GetPositionType.Effect);
            if (!PoolHelper.Instance.CheckPool(poolKey))
            {
                var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
                await UniTask.WaitUntil(() => baseObj.Loaded);
                PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
            }
            await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
            var Eff = PoolHelper.Instance.GetPoolObj(poolKey);
            Eff.transform.SetParent(FightController.Instance.Root_Eff3D);
            Eff.transform.position = new Vector3(transform.position.x, transform.position.y, SceneManager.EffOffZ);

            if (poolKey == "Monster_zd")
            {
                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.FightAtlas, bulletIcon, Eff.GetComponent<SpriteRenderer>()).Forget();
            }

            Eff.SetActive(true);

            Vector2 _direction = (transform.position - tpos);
            _direction.Normalize();
            Eff.transform.rotation = Quaternion.AngleAxis(Mathf.Atan2(_direction.y, _direction.x) * Mathf.Rad2Deg, Vector3.forward);

            Eff.transform.DOMove(tpos, m_MonsterEffTime).SetEase(Ease.Linear).OnComplete(() =>
            {
                Eff.transform.DOKill();
                PoolHelper.Instance.RecycObj(Eff);
            });
        }

        public async UniTask PlayHurt(int skillId)
        {
            string poolKey = "Hurt_1";

            if (skillId > 0)
            {
                if (CFG_Skill.GetById(skillId).SkillProfession == 100)
                    poolKey = "eff_FSAttack";
                else if (CFG_Skill.GetById(skillId).SkillProfession % 10 != 0)
                    poolKey = "Hurt_" + (CFG_Skill.GetById(skillId).SkillProfession % 10);
            }

            AudioManager.Instance.PlaySound((AudioKey.Sound)Enum.Parse(typeof(AudioKey.Sound), poolKey));

            if (!PoolHelper.Instance.CheckPool(poolKey))
            {
                var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
                await UniTask.WaitUntil(() => baseObj.Loaded);
                PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
            }
            await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
            var Eff = PoolHelper.Instance.GetPoolObj(poolKey);
            Eff.transform.SetParent(FightController.Instance.Root_Eff3D);
            Eff.transform.position = new Vector3(transform.position.x, transform.position.y, SceneManager.EffOffZ);
            var skeleton = Eff.GetComponentInChildren<SkeletonAnimation>();
            skeleton.loop = false;
            Eff.SetActive(true);
            skeleton.AnimationState.SetAnimation(0, "hit", false);

            float waitTime = skeleton.AnimationState.GetCurrent(0).Animation.Duration;
            CoroutineCall(waitTime, () =>
            {
                if (Eff != null)
                    PoolHelper.Instance.RecycObj(Eff);
            });
        }

        public void AddBuffers(int[] buffers)
        {
            foreach (var buffer in buffers)
                Debug.Log("---- buffId: " + buffer);

            for (int i = 0; i < buffers.Length; i++)
            {
                if (!m_BufferList.Exists(x => x == buffers[i]))
                    AddOrRemoveBuffer(buffers[i]).Forget();
            }

            CoroutineCall(0.8f, ClearBuff);
            CoroutineCall(m_BuffTime, ClearBuffObjDic);
        }

        public async UniTask AddOrRemoveBuffer(int bufferId, bool isAdd = true)
        {
            if (isAdd)
            {
                if (!m_BufferList.Contains(bufferId))
                    m_BufferList.Add(bufferId);
            }
            else
            {
                if (m_BufferList.Contains(bufferId))
                    m_BufferList.Remove(bufferId);
            }

            string effKey = CFG_BuffEffect.GetById(bufferId).BuffEffect;
            if (!string.IsNullOrEmpty(effKey))
            {
                if (isAdd)
                {
                    string poolKey = effKey;
                    if (!PoolHelper.Instance.CheckPool(poolKey))
                    {
                        var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
                        await UniTask.WaitUntil(() => baseObj.Loaded);
                        PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
                    }
                    await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
                    GameObject effObj = PoolHelper.Instance.GetPoolObj(effKey);
                    effObj.transform.SetParent(transform);
                    effObj.transform.localPosition = Vector3.zero;
                    effObj.gameObject.name = effKey;
                    effObj.SetActive(true);
                    AddJaction(effKey, effObj);
                }
                else
                {
                    ClearOneBuffObj(effKey);
                }
            }
        }

        // 移除buff, 忽略掉持续生效的buff
        public void RemoveUnusedBuffers()
        {
            if (m_BufferList.Count == 0) return;
            var buffs = TurnBasedM.Instance.GetBuffsByUid(TurnBasedFightCtrl.Instance.attackNum, uid);
            List<int> buffers = new();
            foreach (var id in m_BufferList)
            {
                if (!buffs.Exists(x => x.Id == id))
                {
                    buffers.Add(id);
                }
            }

            for (int i = 0; i < buffers.Count; i++)
            {
                AddOrRemoveBuffer(buffers[i], false).Forget();
            }
        }

        public void ClearBuff()
        {
            CoroutineCall(0.3f, RemoveUnusedBuffers);
        }

        public void AddJaction(string key, GameObject jac)
        {
            if (!m_BuffObjList.ContainsKey(key))
            {
                m_BuffObjList.Add(key, jac);
            }
            else
            {
                ClearOneBuffObj(key);
                m_BuffObjList[key] = jac;
            }
        }

        public void ClearBuffObjDic()
        {
            foreach (var item in m_BuffObjList.Values)
            {
                PoolHelper.Instance.RecycObj(item);
            }
            m_BuffObjList.Clear();
        }

        public void ClearOneBuffObj(string key)
        {
            if (m_BuffObjList.ContainsKey(key))
            {
                PoolHelper.Instance.RecycObj(m_BuffObjList[key]);
                m_BuffObjList.Remove(key);
            }
        }

        public void StopAttack()
        {
            StopAllTween();
            ClearCoroutine();
        }

        private void ClearCoroutine()
        {
            if (CoroutineMgr.Instance == null)
                return;

            foreach (var item in coroutines)
            {
                CoroutineMgr.Instance.StopCoroutine(item);
            }
            coroutines.Clear();
        }

        private HashSet<Coroutine> coroutines = new HashSet<Coroutine>();

        private Coroutine CoroutineCall(float time, Action callBack)
        {
            var coroutineItem = CoroutineMgr.Instance.StartCoroutine(IECoroutineCall(time, callBack));
            coroutines.Add(coroutineItem);
            return coroutineItem;
        }

        private System.Collections.IEnumerator IECoroutineCall(float time, Action callBack)
        {
            yield return new WaitForSeconds(time);
            callBack?.Invoke();
        }

        /// <summary>
        /// 受击对手中，第一个活着的对手的uid
        /// </summary>
        /// <returns></returns>
        private long GetFirstLivingUid()
        {
            if (campType == TurnBasedBattle.ActionCamp.BLUE_CAMP)
            {
                List<TurnBasedBattle.TBObject> roleData = new List<TurnBasedBattle.TBObject>();
                if (TurnBasedFightCtrl.Instance.attackNum > 1)
                {
                    var data = TurnBasedFightCtrl.Instance.GetLastRoundData();
                    roleData = data.RedStats;
                }
                else
                {
                    roleData = TurnBasedM.Instance.turnBasedBattleData.OriginRed.Objs;
                }

                foreach (var item in roleData)
                {
                    if (item.Uid != 0 && item.State == TurnBasedBattle.TBObject.ObState.NORMAL)
                        return item.Uid;
                }
                return 0;
            }

            else if (campType == TurnBasedBattle.ActionCamp.RED_CAMP)
            {
                List<TurnBasedBattle.TBObject> roleObjs = new List<TurnBasedBattle.TBObject>();
                if (TurnBasedFightCtrl.Instance.attackNum > 1)
                {
                    var data = TurnBasedFightCtrl.Instance.GetLastRoundData();
                    roleObjs = data.BlueStats;
                }
                else
                    roleObjs = TurnBasedM.Instance.turnBasedBattleData.OriginBlue.Objs;

                foreach (var item in roleObjs)
                {
                    if (item.Uid != 0 && item.State == TurnBasedBattle.TBObject.ObState.NORMAL)
                        return item.Uid;
                }
                return 0;
            }
            return 0;
        }

        public void StopAllTween()
        {
            m_Tween?.Kill();
        }

        public void Dead()
        {
            isDead = true;
            if (m_IsMoved)
            {
                MoveBackToStandPosition(() =>
                {
                    ChangeAnimationState(GameConst.ModelAnim.dead);
                });
            }
            else
            {
                ChangeAnimationState(GameConst.ModelAnim.dead);
            }

            ClearBuffObjDic();
        }

        /// <summary>
        /// 复活
        /// </summary>
        public async UniTask Revive()
        {
            isDead = false;
            ChangeAnimationState(GameConst.ModelAnim.stand);

            string poolKey = "eff_HuiGuangFanZao";
            if (!PoolHelper.Instance.CheckPool(poolKey))
            {
                var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
                await UniTask.WaitUntil(() => baseObj.Loaded);
                PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
            }
            await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
            var Eff = PoolHelper.Instance.GetPoolObj(poolKey);
            Eff.transform.SetParent(TurnBasedFightCtrl.Instance.EffectRoot);
            Eff.transform.position = new Vector3(transform.position.x, transform.position.y, SceneManager.EffOffZ);
            Eff.SetActive(true);
            CoroutineCall(1, () =>
            {
                if (Eff != null)
                    PoolHelper.Instance.RecycObj(Eff);
            });
        }

        public void AddHp(int hp)
        {
            currentHp += hp;
            if (currentHp > maxHp)
                currentHp = maxHp;
        }

        public void RemoveHp(int hp)
        {
            currentHp -= hp;
            if (currentHp < 0)
                currentHp = 0;
        }

        void OnDestroy()
        {
            m_BufferList.Clear();
            StopAllTween();
            ClearCoroutine();
            ClearBuffObjDic();
        }
    }
}