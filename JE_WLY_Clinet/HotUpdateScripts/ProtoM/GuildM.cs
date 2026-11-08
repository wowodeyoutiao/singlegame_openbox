using JEngine.Helper;
using JEngine.Singleton;
using MetaFramework.UI;
using pb;
using System.Collections.Generic;

namespace HotUpdateScripts.ProtoM
{
    public class GuildM : SingletonTemplate<GuildM>
    {
        /// <summary>
        /// 行会列表每页显示的行会数量。
        /// （行会数量不会超过这个数，默认全部显示。）
        /// </summary>
        public const int PageSize = 1000;

        /// <summary>
        /// 工会信息
        /// </summary>
        public GuildInfo guildInfo;

        /// <summary>
        /// 我的职位：1 会长 2 副会长 3 成员
        /// </summary>
        public int roleJob;

        /// <summary>
        /// 我的捐献次数
        /// </summary>
        public int donateNum;

        /// <summary>
        /// 工会成员
        /// </summary>
        public List<GuildMemberInfo> membersInfo;

        /// <summary>
        /// 工会申请列表
        /// </summary>
        /// 
        public List<GuildApplyInfo> applyList;

        /// <summary>
        /// 工会商店购买信息
        /// </summary>
        public List<GuildShopBuyInfo> buyInfo;


        public bool isCheckedStoreRed;

        public string GuildName
        {
            get
            {
                if (guildInfo != null)
                    return guildInfo.Name;
                return string.Empty;
            }
        }


        /// <summary>
        /// 拉取我的行会的所有的信息。
        /// </summary>
        public void GetAllGuildInfo()
        {
            TcpC.C2S_GuildInfoRequest(RoleM.Instance.GuildId);
            TcpC.C2S_GuildDonateInfoRequest();
            TcpC.C2S_GuildShopBuyInfoRequest();
        }

        /// <summary>
        /// 检查行会功能的红点状态。
        /// </summary>
        public void CheckGuildRed()
        {
            isCheckedStoreRed = false;
            if (RoleM.Instance.IsJoinedGuild())
            {
                GetAllGuildInfo();
                CheckRedIsJoinedGuild();
            }
            else
                CheckRedIsNotJoinGuild();
        }

        /// <summary>
        /// 检查未加入行会时的红点状态
        /// </summary>
        public void CheckRedIsNotJoinGuild()
        {
            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Guild_InGuild.ToString(), RedPointManager.Instance.GetRedNum(RedPointConst.Guild_InGuild.ToString()));
        }

        /// <summary>
        /// 检查加入行会之后的红点状态
        /// </summary>
        public void CheckRedIsJoinedGuild()
        {
            CheckGuildSkillRed();
            //CheckGuildMapRed();
        }

        /// <summary>
        /// 剩余可捐献次数。
        /// </summary>
        /// <returns>剩余可捐献次数</returns>
        public int GetRemaingDonateNum()
        {
            return CFG_XMLConfig.xmlConfig.Guild.donatelimitDailytime - donateNum;
        }

        /// <summary>
        /// 检查行会信息红点
        /// </summary>
        public void CheckGuildInfoRed()
        {
            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Guild_Info.ToString(), GetRemaingDonateNum());
        }

        /// <summary>
        /// 检查行会成员申请按钮红点
        /// </summary>
        public void CheckGuildApplyRed()
        {
            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Guild_Apply.ToString(), GetNeedDisposeApplyCount());
        }

        /// <summary>
        /// 需要处理的成员申请数量。
        /// 会长副会长有权处理，成员返回 0。
        /// </summary>
        /// <returns></returns>
        public int GetNeedDisposeApplyCount()
        {
            return roleJob == 1 || roleJob == 2 ? applyList.Count : 0;
        }

        /// <summary>
        /// 检查行会技能按钮红点
        /// </summary>
        public void CheckGuildSkillRed()
        {
            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Guild_Skill.ToString(), (int)RoleM.Instance.GuildWish);
        }

        /// <summary>
        /// 检查行会商店按钮红点
        /// </summary>
        public void CheckGuildStoreRed()
        {
            int num = RedPointManager.Instance.GetRedNum(RedPointConst.Guild_Store.ToString());
            if (num == 1)
            {
                num = IsGuildStoreCanBuy() ? 1 : 0;
            }
            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Guild_Store.ToString(), num);
        }

        /// <summary>
        /// 行会商店可兑换道具的数量。
        /// </summary>
        /// <returns></returns>
        public bool IsGuildStoreCanBuy()
        {
            var storeItems = CFG_GuildStore.GetList();
            foreach (var item in storeItems)
            {
                bool islocked = guildInfo.Level < item.Unlocklevel;
                bool isEnoughBuy = GameUtil.Instance.ComCoinToCheck(item.Price[0].id, item.Price[0].num, null, false);
                if (!islocked && isEnoughBuy)
                    return true;
            }
            return false;
        }

        /// <summary>
        /// 检查行会地图红点
        /// </summary>
        public void CheckGuildMapRed()
        {
            RedPointManager.Instance.RPS.Dispatch(RedPointConst.Guild_Map.ToString(), IsGuildCopyOpen() ? 1 : 0);
        }

        /// <summary>
        /// 行会副本是否开启
        /// </summary>
        /// <returns></returns>
        public bool IsGuildCopyOpen()
        {
            return GameUtil.Instance.CheakCurTime(CFG_XMLConfig.xmlConfig.Guild.guildcopy.starttime, CFG_XMLConfig.xmlConfig.Guild.guildcopy.endtime);
        }

        public bool IsGuildOpen()
        {
            //int lockLv = CFG_XMLConfig.xmlConfig.SwitchOpenLvDic[CFG_XMLConfig.XML_Config.SwitchOpne.Guild];
            //return RoleM.Instance.Level >= lockLv;

            if (GameUtil.Instance.CheckOpenServerIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Guild, false))
                return false;
            if (GameUtil.Instance.CheakIsLock(CFG_XMLConfig.XML_Config.SwitchOpne.Guild, false))
                return false;

            return true;
        }

    }
}
