using JEngine.Helper;
using JEngine.Singleton;
using pb;
using System;
using System.Collections.Generic;

namespace HotUpdateScripts.ProtoM
{
    /// <summary>
    /// 称号信息
    /// </summary>
    public class TitleM : SingletonTemplate<TitleM>
    {
        //称号信息
        private Dictionary<int, TitleInfo> TitleDic = new Dictionary<int, TitleInfo>();

        private Dictionary<int, int> TiltleDayDic = new Dictionary<int, int>();

        /// <summary>
        /// 设置称号信息
        /// </summary>
        /// <param name="_partners"></param>
        public void SetTitles(List<TitleInfo> Titles)
        {
            TitleDic.Clear();
            TiltleDayDic.Clear();
            DateTime date1 = DateTimeOffset.FromUnixTimeSeconds(GameUtil.Instance.GetCurrentTimeSeconds()).LocalDateTime.Date;
            for (int i = 0; i < Titles.Count; i++)
            {
                TitleDic.Add(Titles[i].CfgId, Titles[i]);

                int usedDay = CFG_Title.GetById(Titles[i].CfgId).Period;
                if (usedDay == 0)
                    TiltleDayDic.Add(Titles[i].CfgId, -1);
                else
                {
                    DateTime data2 = DateTimeOffset.FromUnixTimeSeconds(Titles[i].ActivationTime).LocalDateTime.Date;
                    usedDay = usedDay - (date1 - data2).Days;
                    TiltleDayDic.Add(Titles[i].CfgId, usedDay);
                }
            }
        }

        public bool CheckTitleIsGet(int cfgId)
        {
            return TitleDic.ContainsKey(cfgId);
        }

        public int GetTitleStar(int cfgId)
        {
            if (TitleDic.ContainsKey(cfgId))
            {
                return TitleDic[cfgId].Star;
            }
            return 0;
        }

        public int GetTitleDay(int cfgId)
        {
            if (TiltleDayDic.ContainsKey(cfgId))
            {
                return TiltleDayDic[cfgId];
            }
            return 0;
        }

        public void SetTitleStar(int cfgId, int Star)
        {
            TitleDic[cfgId].Star = Star;
        }

        public int GetTitleEndTime(int cfgId)
        {
            if (TitleDic.ContainsKey(cfgId))
            {
                return TitleDic[cfgId].ActivationTime;
            }
            return -1;
        }

        public bool CheckTitleIsEquip(int cfgId, int soltIDx)
        {
            if (soltIDx == 0)
                return RoleM.Instance.TitleId == cfgId;
            else
            {
                return PartnerM.Instance.GetOnePartnerSlotInfo(soltIDx).TitleId == cfgId;
            }
        }

        public int GetTitleBySolt(int soltIDx)
        {
            if (soltIDx == 0)
                return RoleM.Instance.TitleId;
            else
            {
                return PartnerM.Instance.GetOnePartnerSlotInfo(soltIDx).TitleId;
            }
        }



    }

}

