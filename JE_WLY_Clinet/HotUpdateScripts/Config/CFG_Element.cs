using JEngine.Helper;

using System.Collections.Generic;

namespace HotUpdateScripts.Config
{
    public class CFGItem_Element : CFGItemBase
    {
        public int UID;

        public override long GetIdx()
        {
            return UID;
        }

        public int Part;

        public string Partname;

        public int Level;

        public int Rate;

        public AttributeItem[] Levelvalue;

        public Item[] Levelconsume;

        public decimal Breaklevel;

        public Item[] Breakconsume;

        public AttributeItem[] Breakvalue;

        public struct AttributeItem
        {
            public int key;
            public int value;
        }

        public struct Item
        {
            public int id;
            public int num;
        }
    }

    public class CFG_Element : CFGTableBase<CFGItem_Element>
    {
        private static CFG_Element _AllCFG;

        private static CFG_Element AllCFG
        {
            get
            {
                if (_AllCFG == null)
                    _AllCFG = ExcelDataManager.Instance.GetExcelTable<CFG_Element, CFGItem_Element>();

                return _AllCFG;
            }
        }

        private static Dictionary<int, List<CFGItem_Element>> _dicPartElements;

        public static Dictionary<int, List<CFGItem_Element>> DicPartElements
        {
            get
            {
                if (_dicPartElements == null)
                {
                    _dicPartElements = new Dictionary<int, List<CFGItem_Element>>();
                    foreach (var item in AllCFG.DicList)
                    {
                        if (_dicPartElements.ContainsKey(item.Part))
                        {
                            if (_dicPartElements[item.Part] == null)
                            {
                                _dicPartElements[item.Part] = new List<CFGItem_Element>();
                            }
                            _dicPartElements[item.Part].Add(item);
                        }
                        else
                        {
                            _dicPartElements.Add(item.Part, new List<CFGItem_Element>());
                            _dicPartElements[item.Part].Add(item);
                        }
                    }
                }
                return _dicPartElements;
            }
        }

        public static CFG_Element Get()
        {
            return AllCFG;
        }

        public static CFGItem_Element GetById(int uid)
        {
            return AllCFG.DicDatas[uid];
        }

        public static CFGItem_Element GetByPosAndLevel(int pos, int level)
        {
            foreach (var item in AllCFG.DicList)
            {
                if(item.Part == pos && item.Level == level)
                    return item;
            }
            return null;
        }

        /// <summary>
        /// 根据元素部位读取所有元素
        /// </summary>
        /// <param name="part"></param>
        /// <returns></returns>
        public static List<CFGItem_Element> GetElementsByPart(int part)
        {
            DicPartElements.TryGetValue(part, out var elements);
            return elements;
        }

        /// <summary>
        /// 元素的最大等级
        /// </summary>
        /// <param name="part"></param>
        /// <returns></returns>
        public static int MaxLevelOfPart(int part)
        {
            var elements = GetElementsByPart(part);
            if (elements == null) return 0;
            return elements.Count;
        }
    }
}