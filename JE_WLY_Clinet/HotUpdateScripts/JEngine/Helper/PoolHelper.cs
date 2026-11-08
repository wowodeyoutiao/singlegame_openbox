using Cysharp.Threading.Tasks;
using JEngine.Core;
using System;
using System.Collections.Generic;
using UnityEngine;

namespace JEngine.Helper
{
    class PoolInfo
    {
        public string PoolSign;
        public Queue<GameObject> m_Queue = new Queue<GameObject>();
        public GameObject baseGameObj;
        public int MaxCount = 100;
        public PoolInfo(GameObject _baseGameObj, string sign, bool isRelease = false)
        {
            if (baseGameObj == null)
            {
                //Log.Print("_baseGameObj" + _baseGameObj.gameObject.name);
                baseGameObj = GameObject.Instantiate(_baseGameObj, PoolHelper.Instance.transform);
                baseGameObj.gameObject.SetActive(false);
                baseGameObj.gameObject.name = "Base_" + sign;
                PoolSign = sign;
            }
            OnRefreshUsedTime();
            this.isRelease = isRelease;
        }

        private void OnRefreshUsedTime()
        {
            UsedTime = GameUtil.Instance.GetCurrentTimeSeconds();
        }

        private int ReleaseTime = 60 * 5; //5分钟
        private long UsedTime = 0;
        private bool isRelease = false;
        public GameObject GetObj()
        {
            GameObject _CloneObj;
            if (m_Queue.Count > 0)
            {
                _CloneObj = m_Queue.Dequeue();
            }
            else
            {
                _CloneObj = GameObject.Instantiate(baseGameObj);
                _CloneObj.name = PoolSign;
            }
            OnRefreshUsedTime();
            return _CloneObj;
        }
        public void SetObj(GameObject _obj, Transform parent)
        {
            if (m_Queue.Count >= MaxCount)
            {
                UnityEngine.Object.Destroy(_obj);
                return;
            }
            _obj.transform.SetParent(parent);
            _obj.SetActive(false);
            _obj.transform.localPosition = Vector3.zero;
            m_Queue.Enqueue(_obj);
            OnRefreshUsedTime();
        }

        public bool CheakReleasePool()
        {
            if (isRelease && (UsedTime == 0 || GameUtil.Instance.GetCurrentTimeSeconds() - UsedTime > ReleaseTime))
            {
                Clear();
                return true;
            }
            return false;
        }

        public void Clear()
        {
            while (m_Queue.Count != 0)
            {
                var obj = m_Queue.Dequeue();
                if (obj != null)
                    GameObject.Destroy(obj.gameObject);
            }
            m_Queue = null;
            GameObject.Destroy(baseGameObj);
            PoolSign = null;
        }

    }
    public class PoolHelper : MonoBehaviour
    {
        private static Dictionary<string, PoolInfo> PoolDic = new Dictionary<string, PoolInfo>();
        private static PoolHelper _instance;
        public static PoolHelper Instance
        {
            get
            {
                return _instance;
            }
        }

        private void Awake()
        {
            if (_instance == null)
                _instance = this;

            transform.position = Vector3.one * 10000;
            DontDestroyOnLoad(_instance);
            JEngine.Core.CoroutineMgr.Instance.StartCoroutine(PoolCheak());
        }


        private System.Collections.IEnumerator PoolCheak()
        {
            yield return new WaitForSeconds(30);
            if (this.gameObject != null && !this.gameObject.Equals(null))
            {
                List<string> RemovePools = new List<string>();
                foreach (var item in PoolDic)
                {
                    if (item.Value.CheakReleasePool())
                    {
                        RemovePools.Add(item.Key);
                    }
                }
                for (int i = 0; i < RemovePools.Count; i++)
                {
                    PoolDic.Remove(RemovePools[i]);
                }
                RemovePools.Clear();
                CoroutineMgr.Instance.StartCoroutine(PoolCheak());
            }
        }

        /// <summary>
        /// 检测对象池子是否存在
        /// </summary>
        /// <param name="pooKey"></param>
        /// <returns></returns>
        public bool CheckPool(string pooKey)
        {
            return PoolDic.ContainsKey(pooKey);
        }

        public async void LoadObj(string poolKey, string path, Action<GameObject> action)
        {
            if (!PoolHelper.Instance.CheckPool(poolKey))
            {
                await ResouseHelper.Instance.LoadGameObj(path, go =>
                {
                    PoolHelper.Instance.AddPool(poolKey, go);
                });
            }
            await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
            GameObject Eff = PoolHelper.Instance.GetPoolObj(poolKey);
            action.Invoke(Eff);
        }

        /// <summary>
        /// 从对象池中获取
        /// </summary>
        /// <param name="sign"></param>
        /// <param name="baseObj"></param>
        /// <returns></returns>
        public GameObject GetPoolObj(string sign, GameObject baseObj)
        {
            if (!PoolDic.ContainsKey(sign))
            {
                var poolInfo = new PoolInfo(baseObj, sign);
                PoolDic.Add(sign, poolInfo);
            }
            return PoolDic[sign].GetObj();
        }

        public void AddPool(string sign, GameObject baseObj, bool isRelease = true)
        {
            if (!PoolDic.ContainsKey(sign))
            {
                var poolInfo = new PoolInfo(baseObj, sign, isRelease);
                PoolDic.Add(sign, poolInfo);
            }
        }

        public GameObject GetPoolObj(string sign)
        {
            return PoolDic[sign].GetObj();
        }

        public void RecycObj(GameObject _Obj)
        {
            if (_Obj == null || _Obj.Equals(null))
                return;

            var sign = _Obj.gameObject.name;
            if (!PoolDic.ContainsKey(sign))
            {
                Log.PrintWarning($"Miss Pool {sign}");
                Destroy(_Obj);
                return;
            }
            PoolDic[sign].SetObj(_Obj, transform);
        }

        public void RecycObj(string sign, GameObject _Obj)
        {
            if (_Obj == null || _Obj.Equals(null))
                return;
            if (!PoolDic.ContainsKey(sign))
            {
                Destroy(_Obj);
                return;
            }
            PoolDic[sign].SetObj(_Obj, transform);
        }

        void OnDestroy()
        {
            foreach (var item in PoolDic)
            {
                item.Value.Clear();
            }
            PoolDic.Clear();
            _instance = null;
        }

    }


}
