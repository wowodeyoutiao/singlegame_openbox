using JEngine.Singleton;

namespace HotUpdateScripts.ProtoM
{
    /// <summary>
    /// 每日签到
    /// </summary>
    public class SignM : SingletonTemplate<SignM>
    {

    }

    /// <summary>
    /// 签到列表
    /// </summary>
    public struct HTTP_SignListResponse
    {
        public int code { get; set; }
        public HTTP_SignAwsard data { get; set; }
    }

    public struct HTTP_SignAwsard
    {
        //提取的第几轮
        public HTTP_SignDayAward day_award { get; set; }
        //提取的第几天
        public HTTP_SignDayAward[] accumulate_award { get; set; }
    }

    public struct HTTP_SignDayAward
    {
        //提取的第几轮
        public int round { get; set; }
        //提取的第几天
        public int day { get; set; }
        public bool sign { get; set; }
    }

    //签到添加/领取当日
    public struct HTTP_SignOperateRequest
    {
        public int round { get; set; }
        public int day { get; set; }
    }

    //签到添加/领取当日
    public struct HTTP_SignOperateResponse
    {
        public int code { get; set; }
        //提取的第几天
        public bool data { get; set; }
    }

}

