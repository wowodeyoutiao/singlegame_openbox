using JEngine.Core;
using JEngine.Singleton;

namespace HotUpdateScripts.ProtoM
{
    public class GuidanceM : SingletonTemplate<GuidanceM>
    {
        private int passMaxId; // 已完成引导ID
        private bool isPassGuide; // 完成引导
        public bool isGuiding { get; set; } // 新手引导

        private int[] reconnectSteps;
        private int reconnectStepIndex;

        private CFGItem_Guidance cFGItem_Guidance;

        private int[] serverData = new int[1];

        private int skipRoleLv = 15;

        private string GuideKey = "GuideKey";

        public void OnInit()
        {
            passMaxId = 0;
            isGuiding = false;
            isPassGuide = false;
        }

        // guidance 只有一个元素，就是引导ID
        public void SetGuideTypes(int[] guidance)
        {
            if (RoleM.Instance.Level > 12)
            {
                isPassGuide = true;
                return;
            }

            //if (PlayerPrefs.HasKey(GuideKey))
            //    Debug.LogError(PlayerPrefs.GetInt(GuideKey));

            reconnectStepIndex = 0;
            if (guidance == null || guidance.Length == 0)
            {
                passMaxId = 0;
                Log.Print("服务器记录的引导完成的ID：0");
            }
            else
            {
                passMaxId = guidance[0];
                Log.Print("服务器记录的引导完成的ID：" + passMaxId);
                var passCfg = CFG_Guidance.GetById(passMaxId);
                var nextCfg = CFG_Guidance.GetById(passMaxId + 1);
                if (passCfg == null || nextCfg == null)
                    isPassGuide = true;
                else
                {
                    if (passCfg.ReconnectSteps != null && passCfg.ReconnectSteps.Length > 0) // 有可能需要跳过某些步骤
                    {
                        reconnectSteps = passCfg.ReconnectSteps;
                    }
                    else if (passCfg.ReconnectJumpToStage > 0)
                    {
                        // 查找下一段的第一个
                        passCfg = CFG_Guidance.GetFirstByStage(passCfg.ReconnectJumpToStage);
                        if (passCfg == null)
                            isPassGuide = true;
                        else
                            passMaxId = passCfg.ID - 1;
                    }
                    else
                    {
                        if (nextCfg != null && nextCfg.Stage == passCfg.Stage)
                        {
                            // 查找当前阶段的第一个
                            passCfg = CFG_Guidance.GetFirstByStage(passCfg.Stage);
                            passMaxId = passCfg.ID - 1;
                        }
                    }
                }
            }

            GetGuideConfig();
        }

        public bool CanShowBox()
        {
            return true;
            //if (RoleM.Instance.Level >= skipRoleLv || isPassGuide || !isGuiding || passMaxId >= 3)
            //    return true;

            //return false;
        }

        public bool CheckShowGuide()
        {
            //Log.PrintError("---------------------------- ");
            //isPassGuide = true;
            //isGuiding = false;


            if (RoleM.Instance.Level >= skipRoleLv)
            {
                isPassGuide = true;
                isGuiding = false;
                return false;
            }

            if (isPassGuide)
                return false;

            if (CheckTriggerType())
                return true;

            return false;
        }

        // 检查触发
        private bool CheckTriggerType()
        {
            if (cFGItem_Guidance == null) return false;

            switch ((GuidanceTriggerType)cFGItem_Guidance.TriggerType)
            {
                case GuidanceTriggerType.None:
                    return true;

                case GuidanceTriggerType.Level:
                    return RoleM.Instance.Level >= cFGItem_Guidance.TriggerValue;

                case GuidanceTriggerType.Task:
                    return TaskM.Instance.CheckTaskFinish(cFGItem_Guidance.TriggerValue);

                case GuidanceTriggerType.MainBattle:
                    return RoleM.Instance.MainBattleId >= cFGItem_Guidance.TriggerValue;
            }

            return false;
        }

        public void MoveNext()
        {
            Log.Print($"---------------------------------------  完成第{cFGItem_Guidance.Stage}阶段, passId：{cFGItem_Guidance.ID}");
            if (reconnectSteps == null)
                passMaxId++;
            reconnectStepIndex++;
        }

        public void Jump(int id)
        {
            passMaxId = id - 1;
            reconnectStepIndex = 0;
        }

        public CFGItem_Guidance GetGuideConfig()
        {
            if (isPassGuide)
                return null;

            var cfgId = 0;
            if (reconnectSteps != null && reconnectSteps.Length > 0)
            {
                if (reconnectSteps.Length > reconnectStepIndex)
                    cfgId = reconnectSteps[reconnectStepIndex];
                else
                {
                    passMaxId = reconnectSteps[reconnectSteps.Length - 1];
                    reconnectSteps = null;
                }
            }

            if (cfgId == 0)
                cfgId = passMaxId + 1;

            cFGItem_Guidance = CFG_Guidance.GetById(cfgId);
            isPassGuide = cFGItem_Guidance == null;
            isGuiding = !isPassGuide;
            return cFGItem_Guidance;
        }

        public void SaveGuide(int id)
        {
            serverData[0] = id;
            // PlayerPrefs.SetInt(GuideKey, id);
            Log.Print($"保存引导ID到服务器，passMaxId:{id}");
            TcpC.C2S_GuidanceUpdateRequest(serverData);
        }
    }
}