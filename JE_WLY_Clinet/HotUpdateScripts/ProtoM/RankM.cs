using JEngine.Helper;
using JEngine.Singleton;

namespace HotUpdateScripts.ProtoM
{
    /// <summary>
    /// 排行榜
    /// </summary>
    public class RankM : SingletonTemplate<RankM>
    {

        //限时排行版是否开启
        public bool CheakLimetRankOpne()
        {
            foreach (var item in CFG_XMLConfig.xmlConfig.LimitRankDic)
            {
                if (GetLimitRankTime(item.Key) > 0)
                {
                    return true;
                }
            }
            return false;
        }

        public long GetLimitRankTime(int _type)
        {
            long openZeroTime = GameUtil.Instance.ConvertToMidnightSeconds(GameUtil.Instance.GetServerOpenTimeMillisecond());
            return openZeroTime + CFG_XMLConfig.xmlConfig.LimitRankDic[_type] * 86400 - GameUtil.Instance.GetCurrentTimeSeconds();
        }
    }

    public struct HTTP_RankListRequest
    {
        public int type { get; set; }
        public int all { get; set; }
    }

    /// <summary>
    /// 排行榜列表
    /// </summary>
    public struct HTTP_RankListResponse
    {
        public int code { get; set; }
        public HTTP_RankDetail data { get; set; }
    }

    public struct HTTP_RankDetail
    {
        //总排行榜
        public HTTP_RankUserInfo[] rank_all { get; set; }
        //个人榜单信息
        public HTTP_RankUserInfo rank_own { get; set; }
    }

    public struct HTTP_RankUserInfo
    {
        public long uid { get; set; }
        //排行榜类型
        public int type { get; set; }
        public int job { get; set; } //职业
        public int sex { get; set; } //等级
        //名次
        public int num { get; set; }

        //昵称
        public string name { get; set; }

        //数值
        public long value { get; set; }

        // 服务器名字
        public string server_name { get; set; }

        public int avatar_v2 { get; set; }

        public int frame { get; set; }
    }


}

