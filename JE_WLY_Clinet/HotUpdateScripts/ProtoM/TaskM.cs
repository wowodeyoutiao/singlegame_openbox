using JEngine.Core;
using JEngine.Singleton;
using pb;
using System.Collections.Generic;
using System.Linq;

namespace HotUpdateScripts.ProtoM
{
    /// <summary>
    /// 任务信息
    /// </summary>
    public class TaskM : SingletonTemplate<TaskM>
    {
        private Dictionary<TaskClass, List<TaskInfo>> AllTaskDic = new Dictionary<TaskClass, List<TaskInfo>>();
        public void SetAllTasks(List<TaskInfo> tasks)
        {
            AllTaskDic.Clear();
            foreach (var item in tasks)
            {
                if (!AllTaskDic.ContainsKey(item.Class))
                {
                    AllTaskDic.Add(item.Class, new List<TaskInfo>());
                }
                //UnityEngine.Debug.LogError($"{item.CfgId}, {item.Class}, {item.Status}");
                AllTaskDic[item.Class].Add(item);
            }
        }

        public void SendTaskRed(TaskClass taskClass)
        {
            if (taskClass == TaskClass.INVALID_TaskClass || taskClass == TaskClass.DAILY || taskClass == TaskClass.DAILY_ACTIVITY)
            {
                MetaFramework.UI.RedPointManager.Instance.RPS.Dispatch(TaskClass.DAILY.ToString(), TaskM.Instance.TaskRedNumByTaskClass(TaskClass.DAILY));
            }

            if (taskClass == TaskClass.INVALID_TaskClass || taskClass == TaskClass.WEEKLY || taskClass == TaskClass.WEEKLY_ACTIVITY)
            {
                MetaFramework.UI.RedPointManager.Instance.RPS.Dispatch(TaskClass.WEEKLY.ToString(), TaskM.Instance.TaskRedNumByTaskClass(TaskClass.WEEKLY));
            }

            if (taskClass == TaskClass.INVALID_TaskClass || taskClass == TaskClass.ACHIEVEMENT)
            {
                MetaFramework.UI.RedPointManager.Instance.RPS.Dispatch(TaskClass.ACHIEVEMENT.ToString(), TaskM.Instance.TaskRedNumByTaskClass(TaskClass.ACHIEVEMENT));
            }

            if (taskClass == TaskClass.INVALID_TaskClass || taskClass == TaskClass.ROUTINE)
            {
                //MetaFramework.UI.RedPointManager.Instance.RPS.Dispatch(TaskClass.ROUTINE.ToString(), TaskM.Instance.TaskRedNumByTaskClass(TaskClass.ROUTINE));
            }

            if (taskClass == TaskClass.INVALID_TaskClass || taskClass == TaskClass.VIP)
            {
                MetaFramework.UI.RedPointManager.Instance.RPS.Dispatch(TaskClass.VIP.ToString(), TaskM.Instance.TaskRedNumByTaskClass(TaskClass.VIP));
            }

            if (taskClass == TaskClass.INVALID_TaskClass || taskClass == TaskClass.SEVENDAY || taskClass == TaskClass.SEVENDAY_ACTIVITY)
            {
                var num = TaskM.Instance.TaskRedNumByTaskClass(TaskClass.SEVENDAY);
                if (num <= 0)
                    num = TaskM.Instance.TaskRedNumByTaskClass(TaskClass.SEVENDAY_ACTIVITY);

                MetaFramework.UI.RedPointManager.Instance.RPS.Dispatch(TaskClass.SEVENDAY.ToString(), num);
            }
        }

        public int TaskRedNumByTaskClass(TaskClass tClass)
        {
            TaskClass activeClase = TaskClass.INVALID_TaskClass;
            if (tClass == TaskClass.DAILY)
                activeClase = TaskClass.DAILY_ACTIVITY;
            else if (tClass == TaskClass.WEEKLY)
                activeClase = TaskClass.WEEKLY_ACTIVITY;

            if (AllTaskDic.ContainsKey(tClass))
            {
                for (int i = 0; i < AllTaskDic[tClass].Count; i++)
                {
                    if (AllTaskDic[tClass][i].Status == TaskStatus.COMPLETED)
                        return 1;
                }
            }
            if (AllTaskDic.ContainsKey(activeClase))
            {
                for (int i = 0; i < AllTaskDic[activeClase].Count; i++)
                {
                    if (AllTaskDic[activeClase][i].Status == TaskStatus.COMPLETED)
                        return 1;
                }
            }
            return 0;
        }

        /// <summary>
        /// VIP已解锁但未购买的礼包数量
        /// </summary>
        /// <returns></returns>
        public int VIPCanBuyGiftCount()
        {
            int count = 0;
            List<CFGItem_Store> cfg_StoreList = CFG_Store.GetListByStoreType(CFG_Store.StoreType.VIPGift);

            for (int i = 0; i < RoleM.Instance.Vip; i++)
            {
                var storeInfo = cfg_StoreList[RoleM.Instance.Vip];
                if (StoreM.Instance.GetBuyNumByID(storeInfo.ID) < storeInfo.Limite.limitnum)
                {
                    count++;
                }
            }
            return count;
        }

        /// <summary>
        /// VIP可领取奖励数量
        /// </summary>
        /// <returns></returns>
        //public int VIPCanGetRewardCount()
        //{
        //    int count = 0;
        //    List<TaskInfo> taskInfos = TaskM.Instance.GetTasksByClass(TaskClass.VIP).FindAll(v =>
        //    {
        //        int vipLv = CFG_VIP.IsMaxVipLevel(RoleM.Instance.Vip) ? RoleM.Instance.Vip : RoleM.Instance.Vip + 1;
        //        return CFG_Task.GetById(v.CfgId).Typeson == vipLv;
        //    });

        //    foreach (var item in taskInfos)
        //    {
        //        if(item.Status == TaskStatus.REWARDED)
        //        {
        //            count++;
        //        }
        //    }

        //    return count;
        //}

        //根据任务种类获取任务信息
        public List<TaskInfo> GetTasksByClass(TaskClass tClass)
        {
            //单机服：新号可能尚未收到该类任务，直接索引会抛 KeyNotFoundException 并沿 ILRuntime 深栈炸成 stack overflow
            if (!AllTaskDic.ContainsKey(tClass))
                return new List<TaskInfo>();
            if (tClass == TaskClass.ROUTINE || tClass == TaskClass.DAILY_ACTIVITY || tClass == TaskClass.WEEKLY_ACTIVITY)
                AllTaskDic[tClass].Sort((a, b) => { return a.CfgId.CompareTo(b.CfgId); });
            else
                AllTaskDic[tClass].Sort(ItemCompare);
            return AllTaskDic[tClass];
        }

        public bool CheckTaskFinish(int cfgId)
        {
            var taskCfg = CFG_Task.GetById(cfgId);
            if (taskCfg == null)
                return true;

            var taskInfos = GetTasksByClass((TaskClass)taskCfg.Class);
            if (taskInfos == null || taskInfos.Count == 0)
                return true;

            for (int i = 0; i < taskInfos.Count; i++)
            {
                if (taskInfos[i].CfgId == cfgId)
                    return taskInfos[i].Status == TaskStatus.COMPLETED;
            }

            return false;
        }

        //接任务
        public void AddTask(TaskInfo task)
        {
            if (!AllTaskDic.ContainsKey(task.Class))
            {
                Log.Print("AddTask.Class --" + task.Class);
                AllTaskDic.Add(task.Class, new List<TaskInfo>());
            }
            //UnityEngine.Debug.LogError($"{task.CfgId}, {task.Class}, {task.Status}");
            AllTaskDic[task.Class].Add(task);
        }

        //删任务
        public void RemoveTask(TaskInfo task)
        {
            if (AllTaskDic.ContainsKey(task.Class))
            {
                var whereRemove = AllTaskDic[task.Class].FirstOrDefault(t => t.CfgId == task.CfgId);
                AllTaskDic[task.Class].Remove(whereRemove);
            }

        }

        //更新任务
        public void UpdateTask(TaskInfo task)
        {
            RemoveTask(task);
            AddTask(task);
        }
        public int ItemCompare(TaskInfo x, TaskInfo y)
        {
            // 定义任务状态的优先级顺序
            int statusPriority(TaskInfo task)
            {
                switch (task.Status)
                {
                    case TaskStatus.COMPLETED:
                        return 0;
                    case TaskStatus.IN_PROGRESS:
                        return 1;
                    case TaskStatus.REWARDED:
                        return 2;
                    default:
                        return 3;
                }
            }

            // 比较两个任务的状态优先级
            int statusComparison = statusPriority(x).CompareTo(statusPriority(y));
            if (statusComparison != 0)
            {
                // 如果状态优先级不同，则直接返回状态比较结果
                return statusComparison;
            }
            else
            {
                // 如果状态优先级相同，则比较任务类型
                return CFG_Task.GetById(x.CfgId).Type.CompareTo(CFG_Task.GetById(y.CfgId).Type);
            }
        }

    }

}

