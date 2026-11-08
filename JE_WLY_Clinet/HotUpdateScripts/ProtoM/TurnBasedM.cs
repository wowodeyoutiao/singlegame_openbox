using HotUpdateScripts.Config;

using JEngine.Helper;
using JEngine.Singleton;

using pb;

using System.Collections.Generic;

using UnityEngine;

namespace HotUpdateScripts.ProtoM
{
    public class TurnBasedM : SingletonTemplate<TurnBasedM>
    {
        public const int turnBasedRoomId = 22;   // 回合pk房间id
        public const int wuXingRoomId = 1200; // 五行房间id

        public int roomId { get; set; }

        public Vector3 pvpMapPosition => new Vector3(-0.4f, -5.4f, 0);

        #region pvp回合制数据
        public TurnBasedBattle turnBasedBattleData { get; set; }

        public long[] pvpFormationUids { get; set; }

        /// <summary>
        /// 对手uid
        /// </summary>
        public long opponentUid { get; set; }

        public TurnBasedBattle.TBPlayer GetPlayerByCampType(TurnBasedBattle.ActionCamp actionCamp)
        {
            return actionCamp switch
            {
                TurnBasedBattle.ActionCamp.BLUE_CAMP => turnBasedBattleData.OriginBlue,
                TurnBasedBattle.ActionCamp.RED_CAMP => turnBasedBattleData.OriginRed,
                _ => null,
            };
        }

        public TurnBasedBattle.TBObject GetRoleInfo(long uid)
        {
            foreach (var roleInfo in turnBasedBattleData.OriginBlue.Objs)
            {
                if (roleInfo.Uid == uid)
                    return roleInfo;
            }

            foreach (var roleInfo in turnBasedBattleData.OriginRed.Objs)
            {
                if (roleInfo.Uid == uid)
                    return roleInfo;
            }

            return null;
        }

        /// <summary>
        /// 获取回合内的所有攻击
        /// </summary>
        /// <param name="round">回合数</param>
        /// <returns></returns>
        public List<TurnBasedBattle.Round> GetAttackActionsByRound(int round)
        {
            List<TurnBasedBattle.Round> actions = new List<TurnBasedBattle.Round>();
            foreach (var value in turnBasedBattleData.Rounds)
            {
                if (value.RoundId == round)
                    actions.Add(value);
            }
            return actions;
        }

        /// <summary>
        /// 获取攻击数据吧
        /// </summary>
        /// <param name="attackActionNum">攻击序号</param>
        /// <returns></returns>
        public TurnBasedBattle.Round GetAttackAction(int attackActionNum)
        {
            if (attackActionNum < 1 || attackActionNum > turnBasedBattleData.Rounds.Count)
                return null;

            return turnBasedBattleData.Rounds[attackActionNum - 1];
        }

        public TurnBasedBattle.ActionCamp GetFirstAttackCamp()
        {
            if (turnBasedBattleData.Rounds.Count == 0)
                return TurnBasedBattle.ActionCamp.INVALID_CAMP;

            return turnBasedBattleData.Rounds[0].ActionCamp;
        }

        public List<TurnBasedBattle.TBBuff> GetBuffsByUid(int attackNum, long uid)
        {
            List<TurnBasedBattle.TBBuff> buffs = new List<TurnBasedBattle.TBBuff>();
            if (attackNum >= 0 && attackNum < turnBasedBattleData.Rounds.Count)
            {
                var round = turnBasedBattleData.Rounds[attackNum - 1];
                var role = round.BlueStats.Find(x => x.Uid == uid);
                if (role != null)
                {
                    foreach (var buff in role.Buffs)
                        buffs.Add(buff);
                }
                else
                {
                    role = round.RedStats.Find(x => x.Uid == uid);
                    if (role != null)
                    {
                        foreach (var buff in role.Buffs)
                            buffs.Add(buff);
                    }
                }
            }

            return buffs;
        }

        public TurnBasedBattle.TBResult GetFightResult()
        {
            return turnBasedBattleData.Result;
        }
        #endregion

        #region 五行试炼副本
        public long[] wuXingFormationUids { get; set; }

        /// <summary>
        /// 目前层数
        /// </summary>
        public int elementScore { get; set; }

        /// <summary>
        /// 剩余挑战次数
        /// </summary>
        public int elementTimes { get; set; }

        /// <summary>
        /// 主元素 1-金 2-木 3-水 4-火 5-土
        /// </summary>
        public int elementMain { get; set; }

        /// <summary>
        /// 元素和槽位等级。
        /// 共11个Item，从第2个开始使用。
        /// </summary>
        public List<ElementEmpower> empowers { get; set; }

        /// <summary>
        /// 是否需要扫荡 1-需要 0-不需要
        /// </summary>
        public int needSweep { get; set; }

        public int currentPassedLayer => elementScore == 0 ? 0 : CFG_ElementTrials.GetLayerById(elementScore);

        /// <summary>
        /// 当前挑战的层数，UI显示用
        /// </summary>
        public int currentChallengeLayer { get; set; }

        public string mainElementName => GetElementName(elementMain);

        // 五行功能是否开启
        public bool isOpen => IsOpen(false);

        public bool IsOpen(bool showTip)
        {
            return !GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WuXing, showTip) && !GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.WuXing, showTip);
        }

        /// <summary>
        /// 元素名字。
        /// 金木水火土。
        /// </summary>
        /// <param name="elementId">元素id</param>
        /// <param name="withSquareBrackets">名字两侧是否带方括号</param>
        /// <returns></returns>
        public string GetElementName(int elementId, bool withSquareBrackets = false)
        {
            return elementId switch
            {

                1 => withSquareBrackets ? "[金]" : "金",
                2 => withSquareBrackets ? "[木]" : "木",
                3 => withSquareBrackets ? "[水]" : "水",
                4 => withSquareBrackets ? "[火]" : "火",
                5 => withSquareBrackets ? "[土]" : "土",
                _ => null,
            };
        }

        /// <summary>
        /// 更新元素等级
        /// </summary>
        /// <param name="slotId"></param>
        /// <param name="newLevel"></param>
        public void UpdateElementLevel(int slotId, int newLevel)
        {
            empowers[slotId].Level = newLevel;
        }

        /// <summary>
        /// 更新元素突破等级
        /// </summary>
        /// <param name="slotId"></param>
        /// <param name="newLevel"></param>
        public void UpdateElementBreakLevel(int slotId, int newLevel)
        {
            empowers[slotId].BreakLevel = newLevel;
        }

        public int GetCanBuyChallengeCount_WuXing()
        {
            int storeId = CFG_XMLConfig.xmlConfig.Element.needitem;
            return CFG_Store.GetById(storeId).Limite.limitnum - (int)StoreM.Instance.GetBuyNumByID(storeId);
        }

        #endregion
    }
}