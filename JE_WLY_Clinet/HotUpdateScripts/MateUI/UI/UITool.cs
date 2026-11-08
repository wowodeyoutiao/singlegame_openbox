using HotUpdateScripts.ProtoM;
using JEngine.Event;
using System;
using System.Collections.Generic;
using UnityEngine;

namespace MetaFramework.UI
{
    /// <summary>
    /// 操作UI的帮助类，可以直接通过静态类调用
    /// 并且自动转化了泛型
    /// </summary>
    public static class UITool
    {
        public static T Get<T>() where T : class, IUIBase
        {
            string uiName = typeof(T).Name;

            IUIBase ui = UIManager.Instance.Get(uiName);
            if (ui == null)
                return null;

            return ui as T;
        }

        public static bool Check<T>() where T : class, IUIBase
        {
            string uiName = typeof(T).Name;
            return UIManager.Instance.Check(uiName);
        }

        public static bool IsNull<T>() where T : class, IUIBase
        {
            string uiName = typeof(T).Name;
            return UIManager.Instance.IsNull(uiName);
        }

        public static bool IsOpen<T>() where T : class, IUIBase
        {
            string uiName = typeof(T).Name;
            return UIManager.Instance.IsOpen(uiName);
        }

        public static bool IsOpen(string uiName)
        {
            return UIManager.Instance.IsOpen(uiName);
        }

        public static void OpenStack<T>(Action<IUIBase> onComplete = null, params object[] args) where T : class, IUIBase
        {
            string uiName = typeof(T).Name;

            UIManager.Instance.OpenStack(uiName, onComplete, args);
        }

        public static void CloseStack()
        {
            UIManager.Instance.CloseStack();
        }

        public static void Open<T>(Action<IUIBase> onComplete = null, params object[] args) where T : class, IUIBase
        {
            string uiName = typeof(T).Name;

            UIManager.Instance.OpenNormal(uiName, onComplete, args);
        }

        public static void Open(string uiName, Action<IUIBase> onComplete = null, params object[] args)
        {
            UIManager.Instance.OpenNormal(uiName, onComplete, args);
        }

        public static void Close<T>() where T : IUIBase
        {
            string uiName = typeof(T).Name;

            UIManager.Instance.CloseNormal(uiName);
            bool inHall = IsInHall();
            JEvent.defaultEvent.Post(new GameEvent.BackUIHall
            {
                isInHall = inHall
            });
        }

        public static void RemoveUI<T>() where T : IUIBase
        {
            string uiName = typeof(T).Name;

            UIManager.Instance.RemoveUI(uiName);
        }


        public static void LoadUII<T>(Transform baseRoot, Action<IUIBase> onComplete = null)
        {
            string uiName = typeof(T).Name;
            UIManager.Instance.LoadUI(uiName, baseRoot, onComplete).Forget();
        }

        private static Dictionary<string, Texture2D> DicUrlTexture = new Dictionary<string, Texture2D>();

        /// <summary>
        /// 设置玩家头像
        /// </summary>
        /// <param name="url"></param>
        /// <param name="rawImage"></param>
        //public static void SetHeadImg(string url, RawImage rawImage)
        //{
        //    if (DicUrlTexture.ContainsKey(url))
        //    {
        //        rawImage.texture = DicUrlTexture[url];
        //        return;
        //    }
        //    HttpRequestHelper.Instance.CreateImageRequest(url, (BestHTTP.HTTPRequest request, BestHTTP.HTTPResponse response) =>
        //    {
        //        if (response.IsSuccess)
        //        {
        //            rawImage.texture = response.DataAsTexture2D;
        //            DicUrlTexture.Add(url, response.DataAsTexture2D);
        //        }
        //    });
        //}

        /// <summary>
        /// 公告提示
        /// </summary>
        /// <param name="str"></param>
        public static void ShowTips(string str)
        {
            if (!string.IsNullOrEmpty(str))
            {
                HotUpdateScripts.ProtoM.TipM.Instance.AddTips(str);

                if (!Check<UITips>())
                    Open<UITips>();
                //else
                //    Get<UITips>().OnTipShow();
            }
        }

        /// <summary>
        /// 跑马灯
        /// </summary>
        /// <param name="str"></param>
        public static void RunNoticeTips(int id, List<string> strs)
        {
            HotUpdateScripts.ProtoM.TipM.Instance.AddRunNotice(id, strs);
            if (!Check<UINoticeRun>())
                UITool.Open<UINoticeRun>();
        }

        /// <summary>
        /// 道具快捷使用
        /// </summary>
        /// <param name="str"></param>
        public static void ItemUseTips(pb.ItemInfo item)
        {
            if (item != null)
            {
                HotUpdateScripts.ProtoM.TipM.Instance.AddItemUse(item);

                if (!UIItemUesTips.IsOpen && !GuidanceM.Instance.isGuiding)
                {
                    UIItemUesTips.IsOpen = true;
                    Open<UIItemUesTips>();
                }
            }
        }

        /// <summary>
        /// 是否在大厅界面
        /// </summary>
        /// <returns>bool</returns>
        public static bool IsInHall()
        {
            var normalUIList = UIManager.Instance.GetAllOpenedUIByLayer(UILayer.Normal);
            if (normalUIList.Count > 0) return false;
            var topUIList = UIManager.Instance.GetAllOpenedUIByLayer(UILayer.Top);
            if (topUIList.Count > 0) return false;
            return true;
        }
    }
}
