//
// JUIExtensions.cs
//
// Author:
//       JasonXuDeveloper（傑） <jasonxudeveloper@gmail.com>
//
// Copyright (c) 2022 JEngine
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
using Cysharp.Threading.Tasks;

using JEngine.Helper;

using Spine.Unity;

using System.Collections.Generic;

using UnityEngine;
using UnityEngine.UI;

namespace JEngine.UI
{
    public static class JUIExtensions
    {
        private static Camera m_UICamera;

        public static Camera uiCamera
        {
            get
            {
                if (m_UICamera == null)
                {
                    GameObject _uiCamera = GameObject.Find("UICamera");
                    if (_uiCamera)
                        m_UICamera = _uiCamera.GetComponent<Camera>();
                }
                return m_UICamera;
            }
        }

        private static Canvas m_UICanvas;

        public static Canvas uiCanvas
        {
            get
            {
                if (m_UICanvas == null)
                {
                    GameObject _uiCanvas = GameObject.Find("UICanvas");
                    if (_uiCanvas)
                        m_UICanvas = _uiCanvas.GetComponent<Canvas>();
                }
                return m_UICanvas;
            }
        }


        /// <summary>
        /// Get a color from hex code
        /// </summary>
        /// <param name="hex"></param>
        /// <returns></returns>
        public static Color ToColor(this string hex)
        {
            ColorUtility.TryParseHtmlString(hex, out var c);
            return c;
        }

        private static Dictionary<Transform, Text> te_cache = new Dictionary<Transform, Text>();
        /// <summary>
        /// Get text component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <returns></returns>
        public static Text txt(this Transform x)
        {
            if (te_cache.TryGetValue(x, out var ret))
            {
                return ret;
            }
            ret = x.GetComponent<Text>();
            if (ret == null)
            {
                ret = x.gameObject.AddComponent<Text>();
            }
            te_cache[x] = ret;
            return ret;
        }
        /// <summary>
        /// Get text component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <param name="name"></param>
        /// <returns></returns>
        public static Text txt(this Transform x, string name)
        {
            return x.Find(name).txt();
        }


        private static Dictionary<Transform, Button> bu_cache = new Dictionary<Transform, Button>();
        /// <summary>
        /// Get button component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <returns></returns>
        public static Button btn(this Transform x)
        {
            if (bu_cache.TryGetValue(x, out var ret))
            {
                return ret;
            }
            ret = x.GetComponent<Button>();
            if (ret == null)
            {
                ret = x.gameObject.AddComponent<Button>();
            }
            bu_cache[x] = ret;
            return ret;
        }
        /// <summary>
        /// Get button component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <param name="name"></param>
        /// <returns></returns>
        public static Button btn(this Transform x, string name)
        {
            return x.Find(name).btn();
        }

        private static Dictionary<Transform, ExButton> exbu_cache = new Dictionary<Transform, ExButton>();
        /// <summary>
        /// Get button component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <returns></returns>
        public static ExButton exbtn(this Transform x)
        {
            if (exbu_cache.TryGetValue(x, out var ret))
            {
                return ret;
            }
            ret = x.GetComponent<ExButton>();
            if (ret == null)
            {
                ret = x.gameObject.AddComponent<ExButton>();
            }
            exbu_cache[x] = ret;
            return ret;
        }
        /// <summary>
        /// Get button component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <param name="name"></param>
        /// <returns></returns>
        public static ExButton exbtn(this Transform x, string name)
        {
            return x.Find(name).exbtn();
        }


        private static Dictionary<Transform, Outline> ou_cache = new Dictionary<Transform, Outline>();
        /// <summary>
        /// Get outline component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <returns></returns>
        public static Outline outline(this Transform x)
        {
            if (ou_cache.TryGetValue(x, out var ret))
            {
                return ret;
            }
            ret = x.GetComponent<Outline>();
            if (ret == null)
            {
                ret = x.gameObject.AddComponent<Outline>();
            }
            ou_cache[x] = ret;
            return ret;
        }
        /// <summary>
        /// Get outline component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <param name="name"></param>
        /// <returns></returns>
        public static Outline outline(this Transform x, string name)
        {
            return x.Find(name).outline();
        }

        private static Dictionary<Transform, Shadow> sh_cache = new Dictionary<Transform, Shadow>();
        /// <summary>
        /// Get shadow component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <returns></returns>
        public static Shadow shadow(this Transform x)
        {
            if (sh_cache.TryGetValue(x, out var ret))
            {
                return ret;
            }
            ret = x.GetComponent<Shadow>();
            if (ret == null)
            {
                ret = x.gameObject.AddComponent<Shadow>();
            }
            sh_cache[x] = ret;
            return ret;
        }
        /// <summary>
        /// Get shadow component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <param name="name"></param>
        /// <returns></returns>
        public static Shadow shadow(this Transform x, string name)
        {
            return x.Find(name).shadow();
        }


        private static Dictionary<Transform, Image> im_cache = new Dictionary<Transform, Image>();
        /// <summary>
        /// Get image component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <returns></returns>
        public static Image img(this Transform x)
        {
            if (im_cache.TryGetValue(x, out var ret))
            {
                return ret;
            }
            ret = x.GetComponent<Image>();
            if (ret == null)
            {
                ret = x.gameObject.AddComponent<Image>();
            }
            im_cache[x] = ret;
            return ret;
        }
        /// <summary>
        /// Get image component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <param name="name"></param>
        /// <returns></returns>
        public static Image img(this Transform x, string name)
        {
            return x.Find(name).img();
        }


        private static Dictionary<Transform, RectTransform> re_cache = new Dictionary<Transform, RectTransform>();
        /// <summary>
        /// Get recttransform component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <returns></returns>
        public static RectTransform rect(this Transform x)
        {
            if (re_cache.TryGetValue(x, out var ret))
            {
                return ret;
            }
            ret = x.GetComponent<RectTransform>();
            if (ret == null)
            {
                ret = x.gameObject.AddComponent<RectTransform>();
            }
            re_cache[x] = ret;
            return ret;
        }
        /// <summary>
        /// Get recttransform component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <param name="name"></param>
        /// <returns></returns>
        public static RectTransform rect(this Transform x, string name)
        {
            return x.Find(name).rect();
        }

        private static Dictionary<Transform, RawImage> ra_cache = new Dictionary<Transform, RawImage>();
        /// <summary>
        /// Get rawimage component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <returns></returns>
        public static RawImage rawImg(this Transform x)
        {
            if (ra_cache.TryGetValue(x, out var ret))
            {
                return ret;
            }
            ret = x.GetComponent<RawImage>();
            if (ret == null)
            {
                ret = x.gameObject.AddComponent<RawImage>();
            }
            ra_cache[x] = ret;
            return ret;
        }
        /// <summary>
        /// Get rawimage component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <param name="name"></param>
        /// <returns></returns>
        public static RawImage rawimage(this Transform x, string name)
        {
            return x.Find(name).rawImg();
        }


        private static Dictionary<Transform, Slider> sl_cache = new Dictionary<Transform, Slider>();
        /// <summary>
        /// Get slider component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <returns></returns>
        public static Slider slider(this Transform x)
        {
            if (sl_cache.TryGetValue(x, out var ret))
            {
                return ret;
            }
            ret = x.GetComponent<Slider>();
            if (ret == null)
            {
                ret = x.gameObject.AddComponent<Slider>();
            }
            sl_cache[x] = ret;
            return ret;
        }
        /// <summary>
        /// Get slider component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <param name="name"></param>
        /// <returns></returns>
        public static Slider slider(this Transform x, string name)
        {
            return x.Find(name).slider();
        }
        private static Dictionary<Transform, Toggle> to_cache = new Dictionary<Transform, Toggle>();
        /// <summary>
        /// Get toggle component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <returns></returns>
        public static Toggle toggle(this Transform x)
        {
            if (to_cache.TryGetValue(x, out var ret))
            {
                return ret;
            }
            ret = x.GetComponent<Toggle>();
            if (ret == null)
            {
                ret = x.gameObject.AddComponent<Toggle>();
            }
            to_cache[x] = ret;
            return ret;
        }
        /// <summary>
        /// Get toggle component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <param name="name"></param>
        /// <returns></returns>
        public static Toggle toggle(this Transform x, string name)
        {
            return x.Find(name).toggle();
        }


        private static Dictionary<Transform, Scrollbar> sc_cache = new Dictionary<Transform, Scrollbar>();
        /// <summary>
        /// Get scrollbar component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <returns></returns>
        public static Scrollbar scrollbar(this Transform x)
        {
            if (sc_cache.TryGetValue(x, out var ret))
            {
                return ret;
            }
            ret = x.GetComponent<Scrollbar>();
            if (ret == null)
            {
                ret = x.gameObject.AddComponent<Scrollbar>();
            }
            sc_cache[x] = ret;
            return ret;
        }
        /// <summary>
        /// Get scrollbar component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <param name="name"></param>
        /// <returns></returns>
        public static Scrollbar scrollbar(this Transform x, string name)
        {
            return x.Find(name).scrollbar();
        }


        private static Dictionary<Transform, Dropdown> dr_cache = new Dictionary<Transform, Dropdown>();
        /// <summary>
        /// Get dropdown component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <returns></returns>
        public static Dropdown dropdown(this Transform x)
        {
            if (dr_cache.TryGetValue(x, out var ret))
            {
                return ret;
            }
            ret = x.GetComponent<Dropdown>();
            if (ret == null)
            {
                ret = x.gameObject.AddComponent<Dropdown>();
            }
            dr_cache[x] = ret;
            return ret;
        }
        /// <summary>
        /// Get dropdown component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <param name="name"></param>
        /// <returns></returns>
        public static Dropdown dropdown(this Transform x, string name)
        {
            return x.Find(name).dropdown();
        }


        private static Dictionary<Transform, InputField> in_cache = new Dictionary<Transform, InputField>();
        /// <summary>
        /// Get inputfield component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <returns></returns>
        public static InputField input(this Transform x)
        {
            if (in_cache.TryGetValue(x, out var ret))
            {
                return ret;
            }
            ret = x.GetComponent<InputField>();
            if (ret == null)
            {
                ret = x.gameObject.AddComponent<InputField>();
            }
            in_cache[x] = ret;
            return ret;
        }
        /// <summary>
        /// Get inputfield component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <param name="name"></param>
        /// <returns></returns>
        public static InputField input(this Transform x, string name)
        {
            return x.Find(name).input();
        }


        private static Dictionary<Transform, Canvas> ca_cache = new Dictionary<Transform, Canvas>();
        /// <summary>
        /// Get canvas component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <returns></returns>
        public static Canvas canvas(this Transform x)
        {
            if (ca_cache.TryGetValue(x, out var ret))
            {
                return ret;
            }
            ret = x.GetComponent<Canvas>();
            if (ret == null)
            {
                ret = x.gameObject.AddComponent<Canvas>();
            }
            ca_cache[x] = ret;
            return ret;
        }
        /// <summary>
        /// Get canvas component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <param name="name"></param>
        /// <returns></returns>
        public static Canvas canvas(this Transform x, string name)
        {
            return x.Find(name).canvas();
        }


        private static Dictionary<Transform, CanvasGroup> cagp_cache = new Dictionary<Transform, CanvasGroup>();
        /// <summary>
        /// Get canvas component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <returns></returns>
        public static CanvasGroup canvasGroup(this Transform x)
        {
            if (cagp_cache.TryGetValue(x, out var ret))
            {
                return ret;
            }
            ret = x.GetComponent<CanvasGroup>();
            if (ret == null)
            {
                ret = x.gameObject.AddComponent<CanvasGroup>();
            }
            cagp_cache[x] = ret;
            return ret;
        }
        /// <summary>
        /// Get canvas component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <param name="name"></param>
        /// <returns></returns>
        public static CanvasGroup canvasGroup(this Transform x, string name)
        {
            return x.Find(name).canvasGroup();
        }

        private static Dictionary<Transform, TouchEventCom> touch_cache = new Dictionary<Transform, TouchEventCom>();
        /// <summary>
        /// Get TouchEventCom component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <returns></returns>
        public static TouchEventCom touch(this Transform x)
        {
            if (touch_cache.TryGetValue(x, out var ret))
            {
                return ret;
            }
            ret = x.GetComponent<TouchEventCom>();
            if (ret == null)
            {
                ret = x.gameObject.AddComponent<TouchEventCom>();
            }
            touch_cache[x] = ret;
            return ret;
        }
        /// <summary>
        /// Get touch component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <param name="name"></param>
        /// <returns></returns>
        public static TouchEventCom touch(this Transform x, string name)
        {
            return x.Find(name).touch();
        }

        private static Dictionary<Transform, SkeletonGraphic> skeleton_cache = new Dictionary<Transform, SkeletonGraphic>();
        /// <summary>
        /// Get SkeletonGraphic component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <returns></returns>
        public static SkeletonGraphic skeleton(this Transform x)
        {
            if (skeleton_cache.TryGetValue(x, out var ret))
            {
                return ret;
            }
            ret = x.GetComponent<SkeletonGraphic>();
            if (ret == null)
            {
                ret = x.gameObject.AddComponent<SkeletonGraphic>();
            }
            skeleton_cache[x] = ret;
            return ret;
        }

        /// <summary>
        /// Get skeleton component on the transform
        /// </summary>
        /// <param name="x"></param>
        /// <param name="name"></param>
        /// <returns></returns>
        public static SkeletonGraphic skeleton(this Transform x, string name)
        {
            return x.Find(name).skeleton();
        }

        /// <summary>
        /// UI元素相对于canvas的本地坐标
        /// </summary>
        /// <param name="ui">UI元素</param>
        /// <returns>The local coordinates relative to the Canvas</returns>
        public static Vector2 GetUILocalPositionRelativeCanvas(Transform ui)
        {
            bool result = RectTransformUtility.ScreenPointToLocalPointInRectangle(uiCanvas.transform as RectTransform, uiCamera.WorldToScreenPoint(ui.position), uiCamera, out Vector2 localPos);
            if (result) return localPos;

            return Vector2.zero;
        }

        /// <summary>
        /// 将屏幕坐标转换为UI元素相对于Canvas的本地坐标
        /// </summary>
        /// <param name="screenPoint"></param>
        /// <returns>The local coordinates relative to the Canvas</returns>
        public static Vector2 ScreenPointToLocalPointInRectangle(Vector3 screenPoint)
        {
            bool result = RectTransformUtility.ScreenPointToLocalPointInRectangle(uiCanvas.transform as RectTransform, screenPoint, uiCamera, out Vector2 localPos);
            if (result) return localPos;

            return Vector2.zero;
        }

        const string uiGrayMatPath = "Common/Material/UIGray";

        public static void SetGray(this Image x, bool isGray = true)
        {
            if (Object.Equals(x, null)) return;

            if (isGray)
            {
                ResouseHelper.Instance.LoadMaterialAsync(uiGrayMatPath, (mat) =>
                {
                    x.material = mat;
                }).Forget();
            }
            else
            {
                x.material = null;
            }
        }

        public static void SetGray(this Text x, bool isGray = true)
        {
            if (Object.Equals(x, null)) return;

            if (isGray)
            {
                ResouseHelper.Instance.LoadMaterialAsync(uiGrayMatPath, (mat) =>
                {
                    x.material = mat;
                }).Forget();
            }
            else
            {
                x.material = null;
            }
        }
    }
}
