using Cysharp.Threading.Tasks;
using JEngine.LitJson;
using System;
using System.Collections.Generic;
using System.Reflection;

namespace JEngine.Helper
{
    public class ExcelDataManager
    {
        private static ExcelDataManager _instance;
        public static ExcelDataManager Instance
        {
            get
            {
                if (_instance == null)
                {
                    _instance = new ExcelDataManager();
                    _instance.IsLoad = false;
                }
                return _instance;
            }
        }

        public bool IsLoad = false;
        public Dictionary<Type, object> dic = new Dictionary<Type, object>();

        public void Clear()
        {
            dic.Clear();
        }

        public ExcelDataManager()
        {
            //JsonMapper.RegisterExporter((float val, JsonWriter jw) =>
            //{
            //    jw.Write(double.Parse(val.ToString()));
            //});
            //JsonMapper.RegisterImporter((double val) =>
            //{
            //    return (float)val;
            //});
        }

        /// <summary>
        /// 获取一张表
        /// </summary>
        /// <typeparam name="K">Table</typeparam>
        /// <typeparam name="V">Item</typeparam>
        /// <returns></returns>
        public K GetExcelTable<K, V>() where K : CFGTableBase<V> where V : CFGItemBase
        {
            Type type = typeof(K);
            if (dic.ContainsKey(type) && dic[type] is K)
                return dic[type] as K;

            //var jsonName = typeof(V).Name.Replace("CFGItem_", "").ToLower();
            //var json =  ResouseHelper.Instance.LoadTextAsset($"Other/TextAssets/{jsonName}.json");
            //if (json == null) return null;
            //string resultString = $"{{\"jsonDatas\":{json.text}}}";
            //var asset = JsonMapper.ToObject<K>(resultString);
            //if (asset != null) dic.Add(type, asset);

            var asset = JsonMapper.ToObject<K>(jsonDic[type]);
            //if (asset != null) dic.Add(type, asset);
            jsonDic.Remove(type); // 移除已加载的json数据，避免重复加载
            return asset;
        }


        public Dictionary<Type, string> jsonDic = new Dictionary<Type, string>();
        /// <summary>
        /// 预加载所有表
        /// </summary>
        public async void LoadAllCFG(Action<int, int> CallBack)
        {
            var CFGList = CFG_All.ConfigTypes;
            jsonDic.Clear();
            IsLoad = false;
            int AllCount = CFGList.Count;
            CallBack?.Invoke(0, AllCount);

            _ = CFG_XMLConfig.InitAsync();

            for (int j = 0; j < AllCount; j++)
            {
                int i = j;
                string jsonName = (CFGList[i]).Name.Replace("CFG_", "").ToLower();
                bool isOneLoad = false;
                UnityEngine.Debug.Log("[CFG] begin " + (i + 1) + "/" + AllCount + " " + jsonName + ".json");
                ResouseHelper.Instance.LoadTextAssetAsync("Other/TextAssets/" + jsonName + ".json", (jsonStr) =>
                {
                    try
                    {
                        string resultString = "{\"jsonDatas\":" + jsonStr + "}";
                        if (!jsonDic.ContainsKey(CFGList[i]))
                        {
                            jsonDic.Add(CFGList[i], resultString);
                        }
                        MethodInfo initMethod = CFGList[i].GetMethod("Get", BindingFlags.Static | BindingFlags.Public);
                        if (initMethod != null)
                        {
                            // 调用静态方法
                            initMethod.Invoke(null, null);
                        }
                        UnityEngine.Debug.Log("[CFG] done " + (i + 1) + "/" + AllCount + " " + jsonName + ".json");
                        if (i == CFGList.Count - 1)
                        {
                            IsLoad = true;
                            jsonDic.Clear();
                        }
                        CallBack?.Invoke(i + 1, CFGList.Count);
                    }
                    catch (System.Exception e)
                    {
                        UnityEngine.Debug.LogError("[CFG] parse failed " + jsonName + ".json : " + e);
                        CallBack?.Invoke(i + 1, CFGList.Count);
                    }
                    finally
                    {
                        isOneLoad = true;
                    }
                }).Forget();
                // Guard: never hang forever when a TextAsset fails to load and the callback is lost.
                int guard = 0;
                await UniTask.WaitUntil(() => isOneLoad || ++guard > 300);
                if (!isOneLoad)
                {
                    UnityEngine.Debug.LogError("[CFG] TIMEOUT " + (i + 1) + "/" + AllCount + " " + jsonName + ".json");
                    CallBack?.Invoke(i + 1, CFGList.Count);
                }
            }
        }


        /// <summary>
        /// 获取一个数据
        /// </summary>
        /// <typeparam name="K">Table</typeparam>
        /// <typeparam name="V">Item</typeparam>
        /// <param name="id"></param>
        /// <returns></returns>
        public V GetExcelItem<K, V>(long id) where K : CFGTableBase<V> where V : CFGItemBase
        {
            var excelData = GetExcelTable<K, V>();

            if (excelData == null) return null;

            return excelData.GetExcelItem(id);
        }

    }
}