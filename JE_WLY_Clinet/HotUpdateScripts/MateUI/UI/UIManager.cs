using Cysharp.Threading.Tasks;
using JEngine.Core;
using JEngine.Singleton;
using System;
using System.Collections.Generic;
using System.Linq;
using UnityEngine;

namespace MetaFramework.UI
{
    /// <summary>
    /// UI的命名规则：
    /// Prefab：     UIXXX
    /// Control：    UIXXX
    /// Model：      UIModelXXX
    /// View：       UIViewXXX
    /// </summary>
    public class UIManager : SingletonTemplate<UIManager>
    {
        public const string UI_RES_PATH = "Assets/HotUpdateResources/Main/UI/{0}.prefab";

        /// <summary>
        /// UI根节点
        /// </summary>
        private GameObject uiRootGo;

        /// <summary>
        /// UI缓存
        /// </summary>
        private Dictionary<string, IUIBase> uiCacheDic;

        /// <summary>
        /// UI栈
        /// </summary>
        private Stack<IUIBase> uiStack;

        /// <summary>
        /// 普通UI列表
        /// </summary>
        private List<IUIBase> uiList;

        public UIManager()
        {

            uiCacheDic = new Dictionary<string, IUIBase>();
            uiStack = new Stack<IUIBase>();
            uiList = new List<IUIBase>();
        }


        public IUIBase Get(string uiName)
        {
            IUIBase ui = null;
            if (!uiCacheDic.TryGetValue(uiName, out ui))
            {
                Debug.LogError("can not find the ui from cache:" + uiName);
                return null;
            }

            return ui;
        }

        /// <summary>
        /// 根据UI的层拿取所有已打开的UI
        /// </summary>
        /// <param name="layer"></param>
        /// <returns></returns>
        public List<IUIBase> GetAllOpenedUIByLayer(UILayer layer)
        {
            List<IUIBase> _uiList = new List<IUIBase>();
            for (int i = 0; i < uiList.Count; i++)
            {
                if (uiList[i].GetLayer() == layer)
                    _uiList.Add(uiList[i]);
            }

            return _uiList;
        }

        public bool Check(string uiName)
        {
            if (!uiCacheDic.ContainsKey(uiName))
            {
                return false;
            }
            return uiCacheDic[uiName].uiGo.activeSelf;
        }

        public bool IsNull(string uiName)
        {
            return uiCacheDic.ContainsKey(uiName);
        }

        public bool IsOpen(string uiName)
        {
            return uiList.Find(ui => ui.uiName == uiName) != null;
        }

        /// <summary>
        /// 普通打开一个UI
        /// </summary>
        /// <param name="uiName"></param>
        /// <param name="onComplete"></param>
        /// <returns></returns>
        public void OpenNormal(string uiName, Action<IUIBase> onComplete, params object[] args)
        {
            //打开UI，优先从缓存取
            IUIBase ui = null;
            if (!uiCacheDic.TryGetValue(uiName, out ui))
            {
                _ = LoadUIAsync(uiName, (ui) =>
                {
                    OnOpenUI(ui, onComplete, args);
                });
            }
            else
            {
                OnOpenUI(ui, onComplete, args);
            }

            void OnOpenUI(IUIBase ui, Action<IUIBase> onComplete, params object[] args)
            {
                if (!uiList.Contains(ui))
                    uiList.Add(ui);
                ui.uiGo.SetActive(true);
                int canverlayer;
                SetCanvasLayer(ui.uiGo, ui.GetLayer(), out canverlayer);
                ui._canvasOrder = canverlayer;
                ui.OnLayerChange();
                AudioManager.Instance.PlaySound(AudioKey.Sound.UI_OpenAndCloseUI);
                ui.OnEnter(args);

                onComplete?.Invoke(ui);

                UniGameEvent.Send(new GameEvent.UIPanel_Operate() { uiName = uiName, isOpen = true });
            }
        }

        /// <summary>
        /// 普通关闭一个UI
        /// </summary>
        /// <param name="uiName"></param>
        public void CloseNormal(string uiName)
        {
            if (uiList.Count <= 0) return;

            IUIBase ui = uiList.Find(
                delegate (IUIBase ui)
                {
                    return ui.uiName.Equals(uiName);
                });

            if (ui == null)
            {
                Log.PrintWarning("can not find the ui:" + uiName);
                return;
            }

            //if (!ui.i)
            //{
            AudioManager.Instance.PlaySound(AudioKey.Sound.UI_OpenAndCloseUI);
            uiList.Remove(ui);
            ui.uiGo?.SetActive(false);
            ui.OnExit();

            UniGameEvent.Send(new GameEvent.UIPanel_Operate() { uiName = uiName, isOpen = false });
            //}
        }

        /// <summary>
        /// 移除UI从缓存中移除
        /// </summary>
        /// <param name="uiName"></param>
        public void RemoveUI(string uiName)
        {
            if (uiList.Count > 0)
            {
                IUIBase ui = uiList.Find(
                delegate (IUIBase ui)
                {
                    return ui.uiName.Equals(uiName);
                });
                if (ui != null)
                {
                    uiList.Remove(ui);
                }
                else
                {
                    Log.PrintWarning("can not find the ui:" + uiName);
                }
            }
            if (uiCacheDic.ContainsKey(uiName))
            {
                uiCacheDic[uiName].OnDestroy();
                uiCacheDic[uiName] = null;
                uiCacheDic.Remove(uiName);
            }
        }

        /// <summary>
        /// 通过栈打开一个UI
        /// </summary>
        /// <param name="uiName"></param>
        /// <param name="onComplete"></param>
        public void OpenStack(string uiName, Action<IUIBase> onComplete, params object[] args)
        {
            ////栈顶界面暂停
            //if (uiStack.Count > 0)
            //{
            //    IUIBase topUI = uiStack.Peek();
            //    topUI.OnPause();
            //}

            //打开UI，优先从缓存取
            IUIBase ui = null;
            if (!uiCacheDic.TryGetValue(uiName, out ui))
            {
                LoadUIAsync(uiName, (ui) =>
                {
                    OnOpenUI(ui, onComplete, args);
                }).Forget();
            }
            else
            {
                OnOpenUI(ui, onComplete, args);
            }

            void OnOpenUI(IUIBase ui, Action<IUIBase> onComplete, params object[] args)
            {
                //入栈
                uiStack.Push(ui);
                ui.uiGo.SetActive(true);
                ui.OnEnter(args);
                int canverlayer;
                SetCanvasLayer(ui.uiGo, ui.GetLayer(), out canverlayer);
                ui._canvasOrder = canverlayer;
                ui.OnLayerChange();

                onComplete.Invoke(ui);


            }
        }

        /// <summary>
        /// 关闭栈顶的UI
        /// </summary>
        public void CloseStack()
        {
            if (uiStack.Count <= 0) return;

            IUIBase ui = uiStack.Pop();
            ui.uiGo.SetActive(false);
            ui.OnExit();

            ////栈顶界面恢复
            //if (uiStack.Count > 0)
            //{
            //    IUIBase topUI = uiStack.Peek();
            //    topUI.OnResume();
            //}
        }

        /// <summary>
        /// 关闭所有界面
        /// </summary>
        public void CloseAll()
        {
            while (uiStack.Count > 0)
            {
                CloseStack();
            }

            for (int i = uiList.Count - 1; i >= 0; i--)
            {
                CloseNormal(uiList[i].uiName);
            }
        }

        /// <summary>
        /// 关闭所有界面
        /// </summary>
        public void CloseAllExceptHall()
        {
            while (uiStack.Count > 1)
            {
                CloseStack();
            }

            for (int i = uiList.Count - 1; i >= 1; i--)
            {
                if (uiList[i].uiName == "UIHall")
                    continue;
                CloseNormal(uiList[i].uiName);
            }
        }

        /// <summary>
        /// 清除所有界面，包括缓存
        /// </summary>
        public void Clear()
        {
            //CloseAll();

            //foreach (var uiCache in uiCacheDic.Values)
            //{
            //    GameObject.Destroy(uiCache.uiGo);
            //}

            uiCacheDic.Clear();
            uiList.Clear();
            uiStack.Clear();
        }

        private HashSet<string> LoadUIHashSet = new HashSet<string>();
        private async UniTaskVoid LoadUIAsync(string uiName, Action<IUIBase> onComplete)
        {
            // create ui gameobject

            string uiPath = string.Format(UI_RES_PATH, uiName);

            if (LoadUIHashSet.Contains(uiName))
                return;
            LoadUIHashSet.Add(uiName);

            var prefab = new JPrefab(uiPath, true);
            //await prefab.WaitForAsyncLoading();
            await UniTask.WaitUntil(() => prefab.Loaded);
            // get ui class name
            string uiOriginName = uiName.Replace("UI", "");
            string uiModelName = "UIModel" + uiOriginName;

            // create ui class
            Type uiType = Type.GetType(uiName);
            Type uiModelType = Type.GetType(uiModelName);

            IUIBase ui = Activator.CreateInstance(uiType) as IUIBase;
            IUIModel uiModel = Activator.CreateInstance(uiModelType) as IUIModel;

            if (uiRootGo == null)
            {
                uiRootGo = GameObject.Find("UICanvas");
            }
            var parent = uiRootGo;
            if (parent == null) throw new Exception("can not find the ui layer:" + ui.GetLayer());

            GameObject uiGo = GameObject.Instantiate(prefab.Instance, uiRootGo.transform);

            uiGo.name = uiName;
            uiGo.transform.localPosition = Vector3.zero;
            uiGo.transform.localScale = Vector3.one;

            IUIView uiView = uiGo.GetComponent<IUIView>();
            // init class
            uiView.Init(uiGo);
            ui.uiName = uiName;
            ui.uiGo = uiGo;
            ui.uiModel = uiModel;
            ui.uiView = uiView;
            ui.OnInit();

            uiCacheDic.Add(uiName, ui);
            LoadUIHashSet.Remove(uiName);
            // complete
            onComplete?.Invoke(ui);
            //prefab.Dispose();

            // });
        }

        private void SetCanvasLayer(GameObject ui, UILayer _uiLayer, out int canverlayer)
        {
            if (_uiLayer == UILayer.NoLayer)
            {
                canverlayer = ui.GetComponent<Canvas>().sortingOrder;
                return;
            }
            Canvas _canvas = ui.GetComponent<Canvas>();
            int maxLayer = -1;

            var uiList = uiCacheDic.Values.ToList();
            for (int i = 0; i < uiList.Count; i++)
            {
                if (uiList[i].uiGo != null && !uiList[i].uiGo.Equals(null) && uiList[i].uiGo.activeSelf && uiList[i].uiGo != ui && uiList[i].GetLayer() == _uiLayer)
                {
                    maxLayer = Math.Max(maxLayer, uiList[i]._canvasOrder);
                }
            }
            if (maxLayer == -1)
                canverlayer = ((int)_uiLayer - 1) * 200;
            else
                canverlayer = maxLayer + 10;

            _canvas.sortingOrder = canverlayer;
        }

        /// <summary>
        /// 只用来加载的UI 用于大UI上的多个UI分割处理（不归属于UI控制）
        /// </summary>
        /// <param name="tParent"></param>
        /// <param name="args"></param>
        public async UniTaskVoid LoadUI(string uiName, Transform tParent, Action<IUIBase> onComplete)
        {
            //if (uiCacheDic.ContainsKey(uiName))
            //{
            //    Debug.Log("已经存在" + uiName);
            //    return;
            //}

            string uiPath = string.Format(UI_RES_PATH, uiName);
            var prefab = new JPrefab(uiPath, true);
            await UniTask.WaitUntil(() => prefab.Loaded);
            //await prefab.WaitForAsyncLoading();

            string uiOriginName = uiName.Replace("UI", "");
            string uiModelName = "UIModel" + uiOriginName;

            Type uiType = Type.GetType(uiName);
            Type uiModelType = Type.GetType(uiModelName);

            IUIBase ui = Activator.CreateInstance(uiType) as IUIBase;
            IUIModel uiModel = Activator.CreateInstance(uiModelType) as IUIModel;

            GameObject uiGo = GameObject.Instantiate(prefab.Instance, tParent);

            uiGo.name = uiName;

            ui.uiView = uiGo.GetComponent<IUIView>();
            ui.uiView.Init(uiGo);
            ui.uiName = uiName;
            ui.uiGo = uiGo;
            ui.uiModel = uiModel;
            ui.OnInit();

            uiGo.SetActive(true);
            //uiCacheDic.Add(uiName, ui);

            onComplete?.Invoke(ui);
            //prefab.Dispose();
        }
    }
}