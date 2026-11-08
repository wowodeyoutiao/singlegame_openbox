using JEngine.Singleton;
using pb;
using System.Collections.Generic;

namespace HotUpdateScripts.ProtoM
{
    /// <summary>
    /// 提示
    /// </summary>
    public class TipM : SingletonTemplate<TipM>
    {
        //浮动提示
        private Queue<string> FlyTipsQueue = new Queue<string>();
        private int MaxQueue = 20;
        public void AddTips(string addStr)
        {
            if (FlyTipsQueue.Count > MaxQueue)
            {
                FlyTipsQueue.Dequeue();
            }
            FlyTipsQueue.Enqueue(addStr);
        }

        public string GetTips()
        {
            return FlyTipsQueue.Dequeue();
        }

        public bool CheckTips()
        {
            return FlyTipsQueue.Count > 0;
        }

        public void ClearTips()
        {
            FlyTipsQueue.Clear();
        }

        public class RunNoticeClass
        {
            public int Id;
            public string Str;
        }
        //跑马灯
        private Queue<RunNoticeClass> RunNoticeQueue = new Queue<RunNoticeClass>();
        private int MaxRunNotice = 10;

        public void AddRunNotice(int id, List<string> strs)
        {
            if (RunNoticeQueue.Count > MaxRunNotice)
            {
                RunNoticeQueue.Dequeue();
            }
            RunNoticeQueue.Enqueue(new RunNoticeClass { Id = id, Str = CFG_Lamp.GetStr(id, strs) });
        }

        public RunNoticeClass GetRunNotice()
        {
            return RunNoticeQueue.Dequeue();
        }

        public bool CheckRunNotice()
        {
            return RunNoticeQueue.Count > 0;
        }

        public void ClearRunNotice()
        {
            RunNoticeQueue.Clear();
        }

        //道具使用
        private Queue<ItemInfo> ItemUseQueue = new Queue<ItemInfo>();
        private int MaxItemUseNotice = 10;

        public void AddItemUse(ItemInfo item)
        {
            if (ItemUseQueue.Count > MaxItemUseNotice)
            {
                ItemUseQueue.Dequeue();
            }
            ItemUseQueue.Enqueue(item);
        }

        public ItemInfo GetItemUse()
        {
            return ItemUseQueue.Dequeue();
        }

        public bool CheckItemUse()
        {
            return ItemUseQueue.Count > 0;
        }

        public void ClearItemUse()
        {
            ItemUseQueue.Clear();
        }
    }
}

