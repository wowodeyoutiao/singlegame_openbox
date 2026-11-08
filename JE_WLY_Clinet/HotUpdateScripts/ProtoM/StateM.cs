using JEngine.Core;

namespace HotUpdateScripts.ProtoM
{
    public class StateM : Singleton<StateM>
    {
        /// <summary>
        /// 是否发送了升级或加速宝箱的请求？
        /// 升级或加速成功后重置为false, 用于防止重复请求。
        /// </summary>
        public bool isSendUpgradeOrSpeedUpBox;

        /// <summary>
        /// 是否有洗练结果
        /// </summary>
        public bool isHasRefineResult;
    }
}
