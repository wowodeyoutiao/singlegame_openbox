using JEngine.Helper;
using JEngine.Singleton;

namespace HotUpdateScripts.ProtoM
{
    /// <summary>
    /// 特权卡
    /// </summary>
    public class PrivilegeM : SingletonTemplate<PrivilegeM>
    {
        private enum PrivilegeType
        {
            Month = 1,
            Lifetime = 2
        }

        private HTTP_PrivilegeDetailDay[] privilegeData;

        public void SetData(HTTP_PrivilegeDetailDay[] data)
        {
            privilegeData = data;

            CheckSkipAD();
        }

        public void CheckSkipAD()
        {
            GameM.Instance.SkipAD = false;
            if (HadPrivilege(PrivilegeType.Month))
            {
                var cfg = CFG_Privilege.GetById((int)PrivilegeType.Month);
                GameM.Instance.SkipAD = cfg.SkipAd == 1;
            }
        }

        public int GetStaminaUp()
        {
            if (privilegeData != null && privilegeData.Length > 0)
            {
                for (int i = 0; i < privilegeData.Length; i++)
                {
                    if (privilegeData[i].pid == 1 && privilegeData[i].buy == 2 && GameUtil.Instance.GetCurrentTimeSeconds() < privilegeData[i].end_time)
                    {
                        var cfg = CFG_Privilege.GetById(1);
                        return cfg.StaminaUp;
                    }
                }
            }
            return 0;
        }

        public bool CheckQuicklyIsiand()
        {
            if (HadPrivilege(PrivilegeType.Lifetime))
            {
                var cfg = CFG_Privilege.GetById((int)PrivilegeType.Lifetime);
                return cfg.QuicklyIsiand == 1;
            }
            return false;
        }

        public bool CheckQuicklyDaily()
        {
            if (HadPrivilege(PrivilegeType.Month))
            {
                var cfg = CFG_Privilege.GetById((int)PrivilegeType.Month);
                return cfg.QuicklyDaily == 1;
            }
            return false;
        }

        public bool CheckAutoBattle()
        {
            if (HadPrivilege(PrivilegeType.Lifetime))
            {
                var cfg = CFG_Privilege.GetById((int)PrivilegeType.Lifetime);
                return cfg.AutoBattle == 1;
            }
            return false;
        }

        private bool HadPrivilege(PrivilegeType ptype)
        {
            return HadPrivilege((int)ptype);
        }

        public bool HadPrivilege(int id)
        {
            if (privilegeData != null && privilegeData.Length > 0)
            {
                for (int i = 0; i < privilegeData.Length; i++)
                {
                    if (privilegeData[i].pid == id && privilegeData[i].buy == 2 && GameUtil.Instance.GetCurrentTimeSeconds() < privilegeData[i].end_time)
                        return true;
                }
            }

            return false;
        }
    }

    /// <summary>
    /// 特权卡列表
    /// </summary>
    public struct HTTP_PrivilegeListResponse
    {
        public int code { get; set; }
        public HTTP_PrivilegeDetailDay[] data { get; set; }
    }

    public struct HTTP_PrivilegeDetailDay
    {
        //内容编号
        public int pid { get; set; }
        // 1-未购买，2-已购买
        public int buy { get; set; }
        //1-待领取，2-已领取
        public int get { get; set; }

        public long end_time { get; set; }
    }

    //特权卡领取
    public struct HTTP_PrivilegeReceiveDayResponse
    {
        public int code { get; set; }
        public string data { get; set; }
    }

    //特权卡领取
    public struct HTTP_PrivilegeReceiveDayRequset
    {
        public int pid { get; set; }
        //提取的第几天
    }

}

