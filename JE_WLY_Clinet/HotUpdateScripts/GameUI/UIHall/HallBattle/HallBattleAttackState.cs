using Cysharp.Threading.Tasks;
using HotUpdateScripts;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using Spine.Unity;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using static GameEvent;

namespace HallBattle
{
    public class HallBattleAttackState : BaseState<UIHall>
    {
        private class ShowSkillInfo
        {
            public CFGItem_Skill skillCfg;
            public float cooldown;
            public float releaseTime;
            public float hitAddNum;
            public float hitAddPer;

            public ShowSkillInfo(CFGItem_Skill cfg, float cooldown, float add, float per)
            {
                skillCfg = cfg;
                this.cooldown = cooldown;
                hitAddNum = add;
                hitAddPer = per;
                //Debug.LogError($"ShowSkillInfo: {cfg.UUID}, {cooldown}, {add}, {per}");
            }

            public void Act()
            {
                releaseTime = Time.time;
            }

            public bool CanRelease()
            {
                return true;
                //return Time.time >= releaseTime + cooldown;
            }

            public long GetHit(float hit)
            {
                return (long)((hit + hitAddNum) * (1 + hitAddPer * 0.0001f));
            }
        }

        private UIViewHall uIViewHall;

        private bool isBgMove = false;

        private Vector3[] battlePoss = new Vector3[5] { new Vector3(-248, 64, 0), new Vector3(-198, 140, 0), new Vector3(-210, 6, 0), new Vector3(-132, 158, 0), new Vector3(-136, -10, 0) };

        private bool isBossLevel = false;
        private int bossIndex;

        public override void Enter()
        {
            uIViewHall = Owner.GetView();

            Owner.SetPetTalkActive(false);

            Clear();

            LoadRoles();
            ChangeBg();
            if (Owner.curMainBattleCfg.flag == (int)MainBattleFlayType.Boss)
                Owner.EnterBossWithoutLoading(() =>
                {
                    MoveToNextLevel();
                });
            else
                MoveToNextLevel();

            UniGameEvent.AddEventListener<GameEvent.HallBattlePassiveMove>(HallBattlePassiveMove);
            UniGameEvent.AddEventListener<GameEvent.HallBattleSettlement>(HallBattleSettlement);
        }

        // 角色跑动到战场边缘，攻击不到怪物，怪物、背景被动移动
        private void HallBattlePassiveMove(HallBattlePassiveMove passiveMove)
        {
            if (Owner == null || Owner.isNullGo)
                return;

            if(hallMonsters != null && hallMonsters.Count > 0)
            {
                for (int i = 0; i < hallMonsters.Count; i++)
                {
                    hallMonsters[i].PassiveMove(passiveMove.passive);
                }
            }

            isBgMove = passiveMove.passive;
        }

        // 进行下一关
        public void MoveToNextLevel()
        {
            //Debug.LogError("MoveToNextLevel");
            HallBattleM.Instance.SendMainBattlePrepareRequest();

            isBgMove = true;
            ChangeBg();

            // 刷新关卡UI
            Owner.SetMainLevelUIActive();
           
            // 英雄恢复血量
            if (Owner.hallRole.IsState(HallRoleState.RestToBattle))
                Owner.hallRole.ChangeRoleState(HallRoleState.MoveToBattlefield);
            else
                Owner.hallRole.ChangeRoleState(HallRoleState.Revive);
            
            // 重置boss血条
            Owner.SetMainBossHpUIActive();

            // 重置武器掉落
            Owner.ResetNormalReward();

            // 加载怪物
            LoadMonster();

            if (Owner.curMainBattleCfg.flag == (int)MainBattleFlayType.Boss)
                AudioManager.Instance.PlayBgMusic(AudioKey.Sound.BGM_Hall_Boss);
            else
                AudioManager.Instance.PlayBgMusic(AudioKey.Sound.BG_HallFight);
        }

        #region 背景
        // 改变背景
        private void ChangeBg()
        {
            isBossLevel = Owner.curMainBattleCfg.flag == (int)MainBattleFlayType.Boss;

            uIViewHall.bgFarNodeTrans.gameObject.SetActive(true);
            uIViewHall.bgMidNodeTrans.gameObject.SetActive(true);
            uIViewHall.bgNearbyNodeTrans.gameObject.SetActive(true);
            uIViewHall.restBgGo.SetActive(false);

            Owner.nearbyHallBg.RefreshBg(isBossLevel);
            //Owner.midHallBg.RefreshBg(isBossLevel);
            Owner.farHallBg.RefreshBg(isBossLevel);
        }

        /// 移动背景
        private void MoveBg()
        {
            Owner.nearbyHallBg.StartBgMove();
            //Owner.midHallBg.StartBgMove();
            Owner.farHallBg.StartBgMove();

            Owner.NormalRewardMoveWithBg(10);
        }

        private void StopBg()
        {
            Owner.nearbyHallBg.StopBgMove();
            //Owner.midHallBg.StopBgMove();
            Owner.farHallBg.StopBgMove();
        }
        #endregion

        #region 英雄站位
        private void LoadRoles()
        {
            // 目前只加载主角
            for (int i = 0; i < uIViewHall.roleTranss.Length; i++)
            {
                uIViewHall.roleTranss[i].gameObject.SetActive(i == 0);
            }

            uIViewHall.roleTranss[0].localPosition = Vector3.zero;
            Owner.hallRole.InitPosition(battlePoss[0] + new Vector3(-200, 0, 0));
        }
        #endregion

        #region Monster（怪物站位）
        private List<HallMonster> hallMonsters = new List<HallMonster>();

        private Vector2 GetMonsterBirthPosition(int totalCount, int index, float scale)
        {
            var x = Screen.width * 0.5f + Random.Range(100, 300) * index + 100;
            if (Owner.curMainBattleCfg.flag == (int)MainBattleFlayType.Boss)
            {
                totalCount = 1;
                if (bossIndex > 0)
                    x = Screen.width * 0.5f + Random.Range(100, 300) * index + 300;
            }

            var spaceY = 120f / totalCount;
            var y = (battlePoss[0].y - 78.5f) * scale;
            if (totalCount > 1)
                y = 110 - Random.Range(spaceY * 0.25f, spaceY * 0.5f) - spaceY * index;

            return new Vector2(x, y);
        }

        private void LoadMonster()
        {
            var scale = 0.3f;
            if (Owner.curMainBattleCfg.flag == (int)MainBattleFlayType.Elite)
                scale = 0.4f;
            else if (Owner.curMainBattleCfg.flag == (int)MainBattleFlayType.Boss)
                scale = 0.525f;

            uIViewHall.monsterTrans.localScale = Vector3.one;// * scale;

            for (int i = 0; i < hallMonsters.Count; i++)
                hallMonsters[i].gameObject.SetActive(false);

            if(Owner.curMainBattleCfg.flag == (int)MainBattleFlayType.Boss)
            {
                var cfg = Owner.curMainBattleCfg.monster[bossIndex];
                Owner.ResetBossHpUIData(cfg[1]);
                LoadMonster(0, cfg, scale);
            }else
            {
                for (int i = 0; i < Owner.curMainBattleCfg.monster.Length; i++)
                {
                    LoadMonster(i, Owner.curMainBattleCfg.monster[i], scale);
                }
            }  
        }

        private void LoadMonster(int index, int[] info, float scale)
        {
            if (index >= hallMonsters.Count)
            {
                var go = new GameObject("Monster");
                go.transform.SetParent(uIViewHall.monsterTrans);
                go.AddComponent<SkeletonGraphic>();
                hallMonsters.Add(go.AddComponent<HallMonster>());
            }

            var hallMonster = hallMonsters[index].GetComponent<HallMonster>();
            hallMonster.InitWithData(info, GetMonsterBirthPosition(Owner.curMainBattleCfg.monster.Length, index, scale), Owner.curMainBattleCfg.flag == (int)MainBattleFlayType.Boss, scale, Owner.isMopUp);
            hallMonster.AddHPHud(Owner.GetView().norMonsterHPHudGo).Forget();
            hallMonster.AddHpChangeTrans(Owner.GetView().hpChangeNodeTrans);

            var cFGItem_Monster = CFG_Monster.GetById(info[0]);
            GameHelper.Instance.SetAsyncObj(cFGItem_Monster.Model, GameConst.GetMonsterModelByID(cFGItem_Monster.Model), (mGO) =>
            {
                if (!Owner.IsState(HallBattleState.Battle))
                {
                    mGO.SetActive(false);
                    PoolHelper.Instance.RecycObj(mGO);
                    return;
                }
                mGO.SetActive(false);

                hallMonster.gameObject.name = $"{cFGItem_Monster.ID}_{index}";
                hallMonster.gameObject.SetActive(true);

                var skeletonGraphic = hallMonster.gameObject.GetComponent<SkeletonGraphic>();
                skeletonGraphic.skeletonDataAsset = mGO.GetComponentInChildren<SkeletonAnimation>().skeletonDataAsset;
                skeletonGraphic.Initialize(true);

                hallMonster.AddTarget(Owner.hallRole);
                Owner.hallRole.AddTarget(hallMonster);

                hallMonster.ChangeMonsterState(HallMonsterState.MoveToBattlefield);
                PoolHelper.Instance.RecycObj(mGO);
            }).Forget();
        }
        #endregion

        public override void Update(float dTime)
        {
            if (isBgMove)
                MoveBg();
            else
                StopBg();
        }

        #region 结算
        private void HallBattleSettlement(HallBattleSettlement args)
        {
            if (Owner == null || Owner.isNullGo)
                return;

            Owner.RoundOver();

            // 怪物赢了，所有怪物等待结算
            if (hallMonsters != null && hallMonsters.Count > 0)
            {
                for (int i = 0; i < hallMonsters.Count; i++)
                {
                    var hallMonster = hallMonsters[i];
                    if(hallMonster.IsAlive())
                    {
                        hallMonster.ChangeMonsterState(HallMonsterState.WaitingSettlement);
                    }
                }
            }

            // boss 关卡一个一个出
            if (Owner.curMainBattleCfg.flag == (int)MainBattleFlayType.Boss)
            {
                bossIndex++;
                if(bossIndex < Owner.curMainBattleCfg.monster.Length)
                {
                    if(!args.IsRoleWin)
                    {
                        WaitSettlement();
                    }
                    else
                    {
                        GameHelper.Instance.SetDeayTimeDo(() =>
                        {
                            MoveToNextLevel();
                        }, 0.5f / Owner.timeScale);
                    }  
                }
                else
                {
                    // 结算
                    if (args.IsRoleWin)
                        Settlement();
                    else
                        WaitSettlement();
                }
            }
            else
            {
                // 结算
                if (args.IsRoleWin)
                    Settlement();
                else
                    WaitSettlement();
            }

            void WaitSettlement()
            {
                Owner.ShowBossMasterStroke();
                GameHelper.Instance.SetDeayTimeDo(() =>
                {
                    Settlement();
                    Owner.hallRole.ChangeRoleState(HallRoleState.Die);
                }, 0.5f);
            }

            void Settlement()
            {
                bossIndex = 0;
                GameHelper.Instance.SetDeayTimeDo(() =>
                {
                    if(args.IsRoleWin)
                    {
                        Owner.Settlement((back) => { MoveToNextLevel(); });
                    }else
                    {
                        if (Owner != null)
                        {
                            Owner.ShowZhuanLoading();

                            GameHelper.Instance.SetDeayTimeDo(() =>
                            {
                                if (Owner != null)
                                {
                                    Owner.ChangeState(HallBattleState.Rest);
                                    Owner.ShowStrengthenUI();
                                }
                            }, 0.6f / Owner.timeScale);
                        }
                    }
                }, 0.5f / Owner.timeScale);
            }
        }
        #endregion

        private void Clear()
        {
            Debug.Log("hall battle attack clear");
            if(hallMonsters != null && hallMonsters.Count > 0)
            {
                for (int i = 0; i < hallMonsters.Count; i++)
                {
                    hallMonsters[i].gameObject.SetActive(false);
                    GameObject.Destroy(hallMonsters[i].gameObject);
                }
                hallMonsters.Clear();
            }

            UniGameEvent.RemoveEventListener<GameEvent.HallBattlePassiveMove>(HallBattlePassiveMove);
            UniGameEvent.RemoveEventListener<GameEvent.HallBattleSettlement>(HallBattleSettlement);
        }

        public override IEnumerator Exit()
        {
            Clear();

            yield return base.Exit();
        }
    }
}