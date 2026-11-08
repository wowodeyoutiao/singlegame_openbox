using JEngine.Helper;

using System.Collections.Generic;

namespace HotUpdateScripts.Config
{
    public class CFGItem_ElementTrials : CFGItemBase
    {
        public int UID;

        public override long GetIdx()
        {
            return UID;
        }

        public int Levellimit;

        public int Layers;

        public MonsterItem[] Monster;

        public int Timelimit;

        public Item[] Settlement;

        public Item[] Firstpassraward;

        public int[] Preview;

        public string Suggestion;

        public string Mark;

        public struct MonsterItem
        {
            public int id;
            public int position;
        }

        public struct Item
        {
            public int id;
            public int num;
        }
    }

    public class CFG_ElementTrials : CFGTableBase<CFGItem_ElementTrials>
    {
        private static CFG_ElementTrials _AllCFG;

        private static CFG_ElementTrials AllCFG
        {
            get
            {
                if (_AllCFG == null)
                    _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_ElementTrials, CFGItem_ElementTrials>();

                return _AllCFG;
            }
        }

        /// <summary>
        /// 根据层数读取
        /// </summary>
        private static Dictionary<int, CFGItem_ElementTrials> _elementTrial;

        public static Dictionary<int, CFGItem_ElementTrials> ElementTrial
        {
            get 
            {
                if (_elementTrial == null)
                {
                    _elementTrial = new Dictionary<int, CFGItem_ElementTrials>();

                    for (int i = 0; i < _AllCFG.DicList.Count; i++)
                    {
                        var cfg = _AllCFG.DicList[i];
                        _elementTrial.Add(cfg.Layers, cfg);
                    }
                }
                return _elementTrial;
            }
        }

        public static CFG_ElementTrials Get()
        {
            return AllCFG;
        }

        public static CFGItem_ElementTrials GetById(int uid)
        {
            return AllCFG.DicDatas[uid];
        }

        // 根据层读取
        public static CFGItem_ElementTrials GetByLayer(int layer)
        {
            ElementTrial.TryGetValue(layer, out var cfg);
            return cfg;
        }

        public static int GetLayerById(int uid)
        {
            return GetById(uid).Layers;
        }

        // 总层数
        public static int GetLayerCount()
        {
            return AllCFG.DicList.Count;
        }
    }
}