using Cysharp.Threading.Tasks;
using DG.Tweening;
using HotUpdateScripts;
using HotUpdateScripts.ProtoM;
using JEngine.Core;
using JEngine.Helper;
using Spine.Unity;
using System.Linq;
using UnityEngine;

namespace HallBattle
{
    public class HallRole : HallBaseRole
    {
        public override void DoAttack(HallMonster target)
        {
            base.DoAttack(target);

            if(attackTimes >= 0)
            {
                UniGameEvent.Send(new GameEvent.HallBattleRoleAttack() { attackTimes = attackTimes });
            }

            if(specialSkillCfgID != 0)
            {
                var cfg = CFG_SkillMainBattle.GetById(specialSkillCfgID);
                Attack(new ShowSkillInfo(cfg, timeScale), true);
                
                UniGameEvent.Send(new GameEvent.HallBallShowSpecialSkillName() { id = specialSkillCfgID });

                specialSkillCfgID = 0;
                return;
            }

            // 查找可以播放的， Combo == 0 的1-1-2,普通的不考虑CD
            if (norRoleSkillInfos == null || norRoleSkillInfos.Count == 0)
                base.PreLoadRoleSkills();

            if(norRoleSkillInfos != null && norRoleSkillInfos.Count > 0)
            {
                ShowSkillInfo showSkillInfo;

                if (norRoleSkillInfos.Count == 1)
                    showSkillInfo = norRoleSkillInfos[0];
                else
                {
                    if (attackRound == 3)
                    {
                        showSkillInfo = norRoleSkillInfos[1];
                        attackRound = 1;
                    }
                    else
                    {
                        showSkillInfo = norRoleSkillInfos[0];
                        attackRound++;
                    }
                }

                Attack(showSkillInfo, false);
            }
        }

        private void Attack(ShowSkillInfo skillInfo, bool isSpecialSkill)
        {
            // 播放技能
            PlayRoleAtkEffect(skillInfo).Forget();

            // 播放群体受击特效
            var hasFireCollider = false;

            if(skillInfo.skillCfg.HitEffect != null && skillInfo.skillCfg.HitEffect.Length > 0)
            {
                var list = skillInfo.skillCfg.HitEffect.Split(',');
                for (int i = 0; i < list.Length; i++)
                {
                    if (skillInfo.skillCfg.HitType[i] == 2)
                        PlayRoleHitEffect(skillInfo, list[i], i).Forget();
                    else
                    {
                        FireBoxCollider2D(skillInfo, list[i], i, isSpecialSkill).Forget();
                        hasFireCollider = true;
                    }
                }
            }

            if (!hasFireCollider)
            {
                // 发射碰撞盒
                FireBoxCollider2D(skillInfo, null, 0, isSpecialSkill).Forget();
            }  
        }

        private void SpecialSkillBuff(ShowSkillInfo skillInfo)
        {
            if (skillInfo == null)
                return;

            if (!IsAlive())
                return;

            // 女神降临。对方受伤害，已方会回血
            if (skillInfo.skillCfg.ID == 17)
            {
                currentHealth += skillInfo.skillCfg.Damage * 0.5f;
                if (currentHealth > maxHealth)
                    currentHealth = maxHealth;
            }
        }

        #region 技能特效
        private async UniTask PlayRoleAtkEffect(ShowSkillInfo skillInfo)
        {
            if (curRoleState != HallRoleState.Attack)
                return;

            skillInfo.Act();
            // 角色动作
            var actSpeed = 1f;
            if (skillInfo.skillCfg.ActSpeed != 0)
                actSpeed = skillInfo.skillCfg.ActSpeed * 0.1f;

            var skillName = skillInfo.skillCfg.SkillName;
            if(skillInfo.skillCfg.Vocation == 5 && skillName.Contains(','))
            {
                var skillNames = skillName.Split(',');
                skillName = skillNames[job - 1];
            }
            SpineUtil.Instance.PlayOnce(spineTrans, skillName, () => { PlayLoopSpine(GameConst.ModelAnim.stand.ToString()); }, actSpeed);

            // 技能
            string poolKey = skillInfo.skillCfg.EffectName;
            if (string.IsNullOrEmpty(poolKey))
                return;

            var effectInfo = poolKey.Split(',');
            poolKey = effectInfo[0];
            if (!PoolHelper.Instance.CheckPool(poolKey))
            {
                var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
                await UniTask.WaitUntil(() => baseObj.Loaded);

                PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
            }
            await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
            var effGo = PoolHelper.Instance.GetPoolObj(poolKey);
            effGo.transform.SetParent(transform.parent);

            var pos = Vector3.zero;
            if (skillInfo.skillCfg.EffectPositionInfo != null && skillInfo.skillCfg.EffectPositionInfo.Length > 1)
            {
                pos.x = skillInfo.skillCfg.EffectPositionInfo[0];
                pos.y = skillInfo.skillCfg.EffectPositionInfo[1];
            }
            effGo.transform.localPosition = transform.localPosition + pos;

            var effectScale = skillInfo.skillCfg.EffectScaleInfo * 0.1f;
            if(effectScale == 0)
                effectScale = 1;

            effGo.transform.localScale = Vector3.one * scale * effectScale;
            effGo.SetActive(false);
            SetEffectSpineTimeScale(effGo);

            GameHelper.Instance.SetDeayTimeDo(() =>
            {
                if (effGo != null && curRoleState == HallRoleState.Attack)
                {
                    SpecialSkillBuff(skillInfo);
                    effGo.SetActive(true);
                    if(effectInfo.Length > 1)
                    {
                        var spine = effGo.GetComponentInChildren<SkeletonGraphic>();
                        SpineUtil.Instance.PlayOnce(spine, effectInfo[1], null, 1);
                    }

                    if (skillInfo.skillCfg.EffectPositionInfo != null && skillInfo.skillCfg.EffectPositionInfo.Length > 3)
                    {
                        var endPos = new Vector3(skillInfo.skillCfg.EffectPositionInfo[2], skillInfo.skillCfg.EffectPositionInfo[3]);
                        var moveTime = skillInfo.skillCfg.EffectPositionInfo[4] * 0.001f;

                        effGo.transform.DOLocalMove(endPos, moveTime / timeScale);
                    }
                }
            }, skillInfo.skillCfg.EffectStartDelay * 0.001f / timeScale);

            GameHelper.Instance.SetDeayTimeDo(() =>
            {
                if (effGo == null || PoolHelper.Instance == null)
                    return;
                PoolHelper.Instance.RecycObj(effGo);
            }, 8 / timeScale);
        }
        #endregion

        #region 群体受击特效
        private async UniTask PlayRoleHitEffect(ShowSkillInfo skillInfo, string hitEffect, int index)
        {
            if (curRoleState != HallRoleState.Attack)
                return;

            // 技能
            string poolKey = hitEffect;
            if (string.IsNullOrEmpty(poolKey))
                return;

            if (!PoolHelper.Instance.CheckPool(poolKey))
            {
                var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
                await UniTask.WaitUntil(() => baseObj.Loaded);

                PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
            }
            await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
            var effGo = PoolHelper.Instance.GetPoolObj(poolKey);
            effGo.transform.SetParent(transform.parent);
            var offset = (job <= 1 ? 0 : (job == 2 ? 20 : 10));
            effGo.transform.localPosition = transform.localPosition + new Vector3(skillInfo.skillCfg.HitPositionInfo[0 + index * 2] + offset, skillInfo.skillCfg.HitPositionInfo[1 + index * 2], 0);

            var effectScale = skillInfo.skillCfg.HitEffectScaleInfo * 0.1f;
            if (effectScale == 0)
                effectScale = 1;
            effGo.transform.localScale = Vector3.one * scale * effectScale;
            effGo.SetActive(false);
            SetEffectSpineTimeScale(effGo);

            GameHelper.Instance.SetDeayTimeDo(() =>
            {
                if (effGo != null && curRoleState == HallRoleState.Attack)
                {
                    effGo.SetActive(true);
                }
            }, skillInfo.skillCfg.HitDelay[index] * 0.001f / timeScale);

            GameHelper.Instance.SetDeayTimeDo(() =>
            {
                if (effGo == null || PoolHelper.Instance == null)
                    return;
                PoolHelper.Instance.RecycObj(effGo);
            }, 5f / timeScale);
        }
        #endregion

        #region 发射碰撞盒
        private async UniTask FireBoxCollider2D(ShowSkillInfo skillInfo, string hitEffect, int index, bool isSpecialSkill)
        {
            if (curRoleState != HallRoleState.Attack)
                return;

            var poolKey = "UIHall_Role_Fire_Collider2D";
            if (!PoolHelper.Instance.CheckPool(poolKey))
            {
                var go = new GameObject(poolKey);
                PoolHelper.Instance.AddPool(poolKey, go);
            }
            await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
            var colliderGo = PoolHelper.Instance.GetPoolObj(poolKey);
            colliderGo.name = poolKey;
            colliderGo.transform.SetParent(transform.parent);
            colliderGo.transform.localPosition = transform.localPosition + new Vector3(skillInfo.skillCfg.ColliderPositionInfo[0], skillInfo.skillCfg.ColliderPositionInfo[1], 0);

            var colliderScale = skillInfo.skillCfg.HitEffectScaleInfo * 0.1f;
            if (colliderScale == 0)
                colliderScale = 1;
            colliderGo.transform.localScale = Vector3.one * colliderScale;
            colliderGo.SetActive(false);

            var hallEffectInfo = colliderGo.GetComponent<HallEffectInfo>();
            if (hallEffectInfo == null)
                hallEffectInfo = colliderGo.AddComponent<HallEffectInfo>();

            hallEffectInfo.DisableCollider();
            var cType = skillInfo.skillCfg.ColliderSize.Length == 1 ? HallEffectColliderType.Circle : HallEffectColliderType.Box;
            hallEffectInfo.OnInit(RoleM.Instance.TotalCombat * skillInfo.skillCfg.Damage * 0.01f, skillInfo.skillCfg.Effectnum, HallEffectType.Role, job);
            hallEffectInfo.SetHPChangeType(isSpecialSkill, skillInfo.skillCfg.Vocation == 5 ? skillInfo.skillCfg.ID : 0);
            hallEffectInfo.SetSize(skillInfo.skillCfg.ColliderSize, cType);
            hallEffectInfo.SetOffset(skillInfo.skillCfg.ColliderOffset, cType);
            hallEffectInfo.SetHitInfo(skillInfo.skillCfg.HitNum, skillInfo.skillCfg.HitTime * 0.001f);
            hallEffectInfo.SetHitEffectInfo(hitEffect, skillInfo.skillCfg.HitPositionInfo);

            GameHelper.Instance.SetDeayTimeDo(() =>
            {
                if (colliderGo != null && curRoleState == HallRoleState.Attack && skillInfo != null)
                {
                    var actTime = skillInfo.skillCfg.ColliderPositionInfo[2];
                    if(skillInfo.skillCfg.ColliderType == 2)
                    {
                        hallEffectInfo.SetSize(0, cType);
                        colliderGo.SetActive(true);

                        hallEffectInfo.SizeZoomTo(1, actTime * 0.001f / timeScale);
                        hallEffectInfo.EnableCollider();
                    }
                    else
                    {
                        hallEffectInfo.EnableCollider();
                        colliderGo.SetActive(true);
                        var targetPos = colliderGo.transform.localPosition;
                        targetPos.x += skillInfo.skillCfg.MaxDistance;
                        colliderGo.transform.DOLocalMove(targetPos, actTime * 0.001f / timeScale).OnComplete(() =>
                        {
                            if (colliderGo != null)
                                PoolHelper.Instance.RecycObj(colliderGo);
                        });
                    }
                       
                }
            }, skillInfo.skillCfg.ColliderDelay * 0.001f / timeScale);
        }
        #endregion
    }
}