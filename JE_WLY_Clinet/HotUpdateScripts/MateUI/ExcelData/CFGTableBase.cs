using System.Collections.Generic;
using System.Linq;
using UnityEngine;

namespace JEngine.Helper
{
    public class CFGTableBase<T> where T : CFGItemBase
    {
        public T[] jsonDatas;

        private Dictionary<long, T> dicDatas;
        public Dictionary<long, T> DicDatas
        {
            get
            {
                int idx = 0;
                try
                {
                    if (dicDatas == null)
                    {
                        dicDatas = new Dictionary<long, T>();

                        for (int i = 0; i < jsonDatas.Length; i++)
                        {
                            idx = i;
                            dicDatas.Add(jsonDatas[i].GetIdx(), jsonDatas[i]);
                        }
                        jsonDatas = null; // ÊÍ·ÅÄÚ´æ
                    }

                }
                catch (System.Exception e)
                {
                    Debug.LogError(idx);
                    Debug.LogError(e);
                }
                return dicDatas;

            }
        }

        private List<T> _dicList;
        public List<T> DicList
        {
            get
            {
                if (_dicList == null)
                {
                    _dicList = DicDatas.Values.ToList();
                }
                return _dicList;
            }
        }



        public T GetExcelItem(long id)
        {
            if (!DicDatas.ContainsKey(id))
            {
                return null;
            }

            return DicDatas[id];
        }
    }
}