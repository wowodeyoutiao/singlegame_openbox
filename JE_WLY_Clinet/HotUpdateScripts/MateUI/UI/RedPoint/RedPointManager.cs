using pb;
using System;
using System.Collections.Generic;
using System.Linq;

namespace MetaFramework.UI
{
    public class RedPointManager
    {
        public RedPointSystems RPS;

        private static RedPointManager _instance;
        public static RedPointManager Instance
        {
            get
            {
                if (_instance == null)
                {
                    _instance = new RedPointManager();
                    _instance.RPS = new RedPointSystems();
                    _instance.InitNode();
                    _instance.RPS.Init(_instance.baseNode);

                }
                return _instance;
            }
        }
        private RedPointNode baseNode;

        private void InitNode()
        {
            baseNode = new RedPointNode(RedPointConst.BaseNode.ToString());
            {
                RedPointNode SetNode = baseNode.AddChildNode(RedPointConst.Set.ToString());
                {
                    SetNode.AddChildNode(RedPointConst.Set_Email.ToString());
                }

                RedPointNode Hall_TopBtnNode = baseNode.AddChildNode(RedPointConst.Hall_TopBtn.ToString());
                {
                    // 活动（七日签到，七日狂欢）
                    RedPointNode ActiveCenterNode = Hall_TopBtnNode.AddChildNode(RedPointConst.ActiveCenter.ToString());
                    {
                        ActiveCenterNode.AddChildNode(RedPointConst.ActiveCenter_Sign7.ToString());
                        ActiveCenterNode.AddChildNode(TaskClass.SEVENDAY.ToString());
                        ActiveCenterNode.AddChildNode(RedPointConst.ActiveCenter_StarWish.ToString());
                        ActiveCenterNode.AddChildNode(RedPointConst.ActiveCenter_CrazyRefine.ToString());
                        ActiveCenterNode.AddChildNode(RedPointConst.ActiveCenter_RefineFund.ToString());
                        ActiveCenterNode.AddChildNode(RedPointConst.ActiveCenter_CallFund.ToString());
                        ActiveCenterNode.AddChildNode(RedPointConst.ActiveCenter_CopyFund.ToString());
                    }

                    //Hall_TopBtnNode.AddChildNode(RedPointConst.SignDay.ToString());
                    Hall_TopBtnNode.AddChildNode(RedPointConst.FirstPay.ToString());

                    //Hall_TopBtnNode.AddChildNode(RedPointConst.WorldBoss.ToString());
                    //Hall_TopBtnNode.AddChildNode(RedPointConst.DayDo.ToString());
                    //Hall_TopBtnNode.AddChildNode(RedPointConst.TurnBox.ToString());
                }

                RedPointNode SummaryNode = baseNode.AddChildNode(RedPointConst.Summary.ToString());
                {
                    RedPointNode RankNode = SummaryNode.AddChildNode(RedPointConst.Rank.ToString());
                    {
                        RankNode.AddChildNode(RedPointConst.Rank_Level.ToString());
                        RankNode.AddChildNode(RedPointConst.Rank_Power.ToString());
                        RankNode.AddChildNode(RedPointConst.Rank_SiLian.ToString());
                        RankNode.AddChildNode(RedPointConst.Rank_ZhanHun.ToString());
                        RankNode.AddChildNode(RedPointConst.Rank_HuiHe.ToString());
                    }

                    // 任务
                    RedPointNode TaskNode = SummaryNode.AddChildNode(RedPointConst.Task.ToString());
                    {
                        TaskNode.AddChildNode(RedPointConst.SignDay.ToString());
                        TaskNode.AddChildNode(pb.TaskClass.DAILY.ToString());
                        //TaskNode.AddChildNode(pb.TaskClass.ROUTINE.ToString());
                        TaskNode.AddChildNode(pb.TaskClass.WEEKLY.ToString());
                        TaskNode.AddChildNode(pb.TaskClass.ACHIEVEMENT.ToString());
                    }

                    // 每日礼包、每日特价、成长基金、宠物礼包、伙伴礼包、月卡
                    RedPointNode GiftBagNode = SummaryNode.AddChildNode(RedPointConst.GiftBag.ToString());
                    {
                        GiftBagNode.AddChildNode(RedPointConst.GiftBag_PetGift.ToString());
                        GiftBagNode.AddChildNode(RedPointConst.GiftBag_PartnerGift.ToString());
                        GiftBagNode.AddChildNode(RedPointConst.GiftBag_Gemstone.ToString());
                        GiftBagNode.AddChildNode(RedPointConst.GiftBag_Element.ToString());
                        GiftBagNode.AddChildNode(RedPointConst.GiftBag_Blessing.ToString());
                        GiftBagNode.AddChildNode(RedPointConst.GiftBag_Refine.ToString());
                    }

                    RedPointNode SpecialOffer = SummaryNode.AddChildNode(RedPointConst.SpecialOffer.ToString());
                    {
                        SpecialOffer.AddChildNode(RedPointConst.GiftBag_DailyGift.ToString());
                        SpecialOffer.AddChildNode(RedPointConst.GiftBag_GiftPay.ToString());
                        SpecialOffer.AddChildNode(RedPointConst.GiftBag_GrowthFund.ToString());
                        SpecialOffer.AddChildNode(RedPointConst.GiftBag_AdGift.ToString());
                        SpecialOffer.AddChildNode(RedPointConst.GiftBag_Power.ToString());
                        SpecialOffer.AddChildNode(RedPointConst.GiftBag_Time.ToString());
                        SpecialOffer.AddChildNode(RedPointConst.PrivilegeCard.ToString());
                    }

                }

                RedPointNode Hall_LeftBtnNode = baseNode.AddChildNode(RedPointConst.Hall_LeftBtn.ToString());
                {
                    Hall_LeftBtnNode.AddChildNode(RedPointConst.LimitedBounty.ToString());
                    Hall_LeftBtnNode.AddChildNode(RedPointConst.RouletteWheel.ToString());
                }


                RedPointNode Hall_VipBtnNode = baseNode.AddChildNode(RedPointConst.VIP_Enter.ToString());
                {
                    //Hall_VipBtnNode.AddChildNode(RedPointConst.VIP_Buy.ToString());
                    Hall_VipBtnNode.AddChildNode(TaskClass.VIP.ToString());
                }

                var BoxLevelNode = baseNode.AddChildNode(RedPointConst.BoxLevel.ToString());
                {
                    BoxLevelNode.AddChildNode(RedPointConst.BoxLevel1.ToString());
                }

                baseNode.AddChildNode(RedPointConst.PVE.ToString());
                baseNode.AddChildNode(RedPointConst.Tower_Get.ToString());

                RedPointNode BagNode = baseNode.AddChildNode(RedPointConst.Bag.ToString());
                {
                    BagNode.AddChildNode(RedPointConst.Bag_Syn.ToString());
                }

                //经验池
                //RedPointNode ExpPool_First = baseNode.AddChildNode(RedPointConst.ExpPool_First.ToString());
                //{
                //    //ExpPool_First.AddChildNode(RedPointConst.ExpPool.ToString());
                //}

                RedPointNode RoleNode = baseNode.AddChildNode(RedPointConst.Role.ToString());
                {
                    //一键装备
                    RoleNode.AddChildNode(RedPointConst.Role_Equip.ToString());

                    //宠物
                    RedPointNode Pet = RoleNode.AddChildNode(RedPointConst.Role_Pet.ToString());
                    {
                        Pet.AddChildNode(RedPointConst.Pet_Up.ToString());
                        Pet.AddChildNode(RedPointConst.Pet_Award.ToString());

                        RedPointNode PetNode = Pet.AddChildNode(RedPointConst.Pet_UpStar.ToString());
                        {
                            PetNode.AddChildNode(RedPointConst.Pet_UpStar_1.ToString());
                        }
                    }

                    //伙伴
                    RoleNode.AddChildNode(RedPointConst.Partner_Into.ToString());
                    RoleNode.AddChildNode(RedPointConst.Partner_UpStar.ToString());


                    RedPointNode Build2 = RoleNode.AddChildNode(RedPointConst.Role_Skill.ToString());
                    {
                        RedPointNode Build3 = Build2.AddChildNode(RedPointConst.Role_Skill_Build.ToString());
                        {

                        }
                    }

                    //打造
                    RedPointNode Build = RoleNode.AddChildNode(RedPointConst.Role_Build.ToString());
                    {
                        RedPointNode EquipNode = Build.AddChildNode(RedPointConst.Equip.ToString());
                        {
                            EquipNode.AddChildNode(RedPointConst.Equip_Reinforce.ToString());
                            EquipNode.AddChildNode(RedPointConst.Equip_Star.ToString());
                            EquipNode.AddChildNode(RedPointConst.Equip_Insert.ToString());
                            EquipNode.AddChildNode(RedPointConst.Equip_Refine.ToString());

                        }

                        RedPointNode SkillNode = Build.AddChildNode(RedPointConst.Skill.ToString());
                        {
                            SkillNode.AddChildNode(RedPointConst.Skill_Advance.ToString());
                            SkillNode.AddChildNode(RedPointConst.Skill_UpLv.ToString());
                            SkillNode.AddChildNode(RedPointConst.Skill_Setting.ToString());
                        }

                    }

                    //称号
                    RedPointNode RoleTitleNode = RoleNode.AddChildNode(RedPointConst.Role_Title.ToString());
                    {
                        RoleTitleNode.AddChildNode(RedPointConst.Role_Title_Active.ToString());
                        RoleTitleNode.AddChildNode(RedPointConst.Role_Title_UpStar.ToString());
                    }
                }

                RedPointNode PartnerCallNode = baseNode.AddChildNode(RedPointConst.Partner_Call.ToString());
                {
                    PartnerCallNode.AddChildNode(RedPointConst.Partner_Call_Free.ToString());
                    PartnerCallNode.AddChildNode(RedPointConst.Role_Bond.ToString());
                }

                RedPointNode GuildNode = baseNode.AddChildNode(RedPointConst.Guild.ToString());
                {
                    RedPointNode GuildTagNode = GuildNode.AddChildNode(RedPointConst.Guild_Tag.ToString());
                    {
                        GuildTagNode.AddChildNode(RedPointConst.Guild_InGuild.ToString());
                        GuildTagNode.AddChildNode(RedPointConst.Guild_Info.ToString());
                        RedPointNode GuildMember = GuildTagNode.AddChildNode(RedPointConst.Guild_Member.ToString());
                        {
                            GuildMember.AddChildNode(RedPointConst.Guild_Apply.ToString());
                        }
                        GuildTagNode.AddChildNode(RedPointConst.Guild_Skill.ToString());
                        GuildTagNode.AddChildNode(RedPointConst.Guild_Store.ToString());
                        GuildTagNode.AddChildNode(RedPointConst.Guild_Map.ToString());
                        //GuildTagNode.AddChildNode(RedPointConst.Guild_List.ToString());

                    }

                    //RedPointNode GuildFriendTagNode = GuildNode.AddChildNode(RedPointConst.GuildFirend_Tag.ToString());

                    //RedPointNode GuildExtendTagNode = GuildNode.AddChildNode(RedPointConst.GuildExtend_Tag.ToString());
                }


                RedPointNode PlayNode = baseNode.AddChildNode(RedPointConst.Play.ToString());
                {
                    RedPointNode ChallengeNode = PlayNode.AddChildNode(RedPointConst.Challenge.ToString());
                    {
                        RedPointNode Challenge_FightNode = ChallengeNode.AddChildNode(RedPointConst.Challenge_Fight.ToString());
                        {
                            Challenge_FightNode.AddChildNode(RedPointConst.Challenge_Fight_Click.ToString());
                            Challenge_FightNode.AddChildNode(RedPointConst.Challenge_Fight_Award.ToString());
                        }

                        ChallengeNode.AddChildNode(RedPointConst.Challenge_Boss.ToString());
                        ChallengeNode.AddChildNode(RedPointConst.Challenge_Relice.ToString());
                        ChallengeNode.AddChildNode(RedPointConst.Challenge_OnHook.ToString());
                        RedPointNode Challenge_TowerNode = ChallengeNode.AddChildNode(RedPointConst.Challenge_Tower.ToString());
                        {
                            //Challenge_TowerNode.AddChildNode(RedPointConst.Challenge_Tower_Award.ToString());
                            Challenge_TowerNode.AddChildNode(RedPointConst.Challenge_Tower_Click.ToString());
                        }
                    }

                    RedPointNode LimetNode = PlayNode.AddChildNode(RedPointConst.Limet.ToString());
                    {
                        LimetNode.AddChildNode(RedPointConst.PVP_More.ToString());
                        LimetNode.AddChildNode(RedPointConst.PVP_One.ToString());
                    }
                    PlayNode.AddChildNode(RedPointConst.DayDo.ToString());
                    PlayNode.AddChildNode(RedPointConst.WuXing.ToString());
                    PlayNode.AddChildNode(RedPointConst.WarFight.ToString());
                }

            }
        }

        //本地存储红点事件
        private Dictionary<string, int> DicLocalRed = new Dictionary<string, int>
        {
            //[RedPointConst.GiftBag_GiftPay.ToString()] = 1,
            //[RedPointConst.GiftBag_PetGift.ToString()] = 1,
            //[RedPointConst.GiftBag_PartnerGift.ToString()] = 1,
            [RedPointConst.Equip_Star.ToString()] = 1,
            [RedPointConst.Equip_Refine.ToString()] = 1,
            [RedPointConst.Skill_Setting.ToString()] = 1,

            [RedPointConst.Rank_Level.ToString()] = 1,
            [RedPointConst.Rank_Power.ToString()] = 1,
            [RedPointConst.Rank_ZhanHun.ToString()] = 1,
            [RedPointConst.Rank_SiLian.ToString()] = 1,
            [RedPointConst.Rank_HuiHe.ToString()] = 1,

            [RedPointConst.Guild_InGuild.ToString()] = 1,
            [RedPointConst.Guild_Store.ToString()] = 1,

            [RedPointConst.Challenge_Fight_Click.ToString()] = 1,
            [RedPointConst.Challenge_Tower_Click.ToString()] = 1,

            [RedPointConst.Bag_Syn.ToString()] = 1,

            [RedPointConst.WarFight.ToString()] = 1,
        };

        public void ClearRed()
        {
            foreach (var item in DicLocalRed.Keys.ToList())
            {
                DicLocalRed[item] = 1;
            }
        }

        public int GetRedNum(string sKey)
        {
            return DicLocalRed[sKey];
        }

        public void SetRedClose(string sKey, Action callBack)
        {
            DicLocalRed[sKey] = 0;
            callBack?.Invoke();
        }
    }
}
