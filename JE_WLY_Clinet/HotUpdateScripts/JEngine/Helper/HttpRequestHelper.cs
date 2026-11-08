using BestHTTP;
using JEngine.Core;
using JEngine.LitJson;
using JEngine.Singleton;
using MetaFramework.UI;
using System;
using System.Collections.Generic;
using UnityEngine;

namespace JEngine.Helper
{
    public class HttpRequestHelper : SingletonTemplate<HttpRequestHelper>
    {
        public delegate void DelegateProgress(long requestID, float progress, float size);
        public delegate void DelegateTextFinish(long requestID, string txt);
        public delegate void DelegateBytesFinish(long requestID, byte[] bytes);
        public delegate void DelegateFileFinish(long requestID);
        public delegate void DelegateError(long requestID, string msg, int code);
        public delegate void DelegateUrlError(long requestID, string msg, int code);
        public delegate void DelegateUrlFinish(long requestID, string msg);

        //struct RequestData
        //{
        //    public long guid;
        //    public bool isTxt;
        //}

        public DelegateProgress OnDownloadProgress;
        public DelegateFileFinish OnDownLoadFileFinish;
        public DelegateTextFinish OnDownLoadTextFinish;
        public DelegateBytesFinish OnDownLoadBytesFinish;
        public DelegateError OnDownloadError;
        public DelegateUrlFinish OnDownLoadUrlFinish;
        public DelegateUrlError OnDownloadUrlError;


        Dictionary<long, HTTPRequest> m_requestSave = new Dictionary<long, HTTPRequest>();

        public void StopRequest(long guid)
        {
            if (m_requestSave.ContainsKey(guid))
            {
                var quest = m_requestSave[guid];
                if (quest != null)
                {
                    quest.Tag = null;
                    if (quest.State < HTTPRequestStates.Finished)
                        quest.Abort();
                }
                m_requestSave.Remove(guid);
            }
        }

        public void StopAllRequest()
        {
            foreach (var kvp in m_requestSave)
            {
                if (kvp.Value != null)
                {
                    kvp.Value.Tag = null;
                    if (kvp.Value.State < HTTPRequestStates.Finished)
                        kvp.Value.Abort();
                }
            }
            m_requestSave.Clear();
        }


        //void OnDownloadFileProgress(HTTPRequest request, long downloaded, long length)
        //{
        //    var data = (RequestData)request.Tag;
        //    if (OnDownloadProgress != null)
        //    {
        //        OnDownloadProgress(data.guid, downloaded, length);
        //    }
        //}



        public System.DateTime ConvertLongToDateTime(long timeStamp)
        {
            System.DateTime dtStart = System.TimeZone.CurrentTimeZone.ToLocalTime(new System.DateTime(1970, 1, 1));
            long lTime = long.Parse(timeStamp + "0000000");
            System.TimeSpan toNow = new System.TimeSpan(lTime);
            return dtStart.Add(toNow);
        }

        public string ConvertDateTimeToLong()
        {
            System.DateTime startTime = System.TimeZone.CurrentTimeZone.ToLocalTime(new System.DateTime(1970, 1, 1));
            long t = (long)(System.DateTime.Now - startTime).TotalSeconds;
            return t.ToString();
        }

        /// <summary>
        /// 单json 传输
        /// </summary>
        /// <typeparam name="T"></typeparam>
        /// <param name="url"></param>
        /// <param name="htpMethod"></param>
        /// <param name="postStr"></param>
        /// <param name="YesFun"></param>
        /// <param name="NoFun"></param>
        //public static void SendPost<T>(string url, HTTPMethods htpMethod, string postStr,Action<T> YesFun = null, Action NoFun = null) where T : struct
        //{
        //    Debug.Log("Post " + postStr);
        //        var htpMa = new HTTPRequest(new Uri(url), htpMethod, ((HTTPRequest httpRequest, HTTPResponse httpResponse) =>
        //        {
        //            if (httpResponse != null)
        //            {
        //                if (httpResponse.IsSuccess)
        //                {
        //                    Debug.Log("httpResponse.DataAsText + " + httpResponse.DataAsText);
        //                    T backData = JsonMapper.ToObject<T>(httpResponse.DataAsText);
        //                    YesFun?.Invoke(backData);
        //                }
        //                else
        //                {
        //                    NoFun?.Invoke();
        //                }
        //            }
        //            else
        //            {
        //                NoFun?.Invoke();
        //            }

        //            Debug.Log("httpResponse.StatusCode  " + httpResponse.StatusCode);
        //        }));
        //        htpMa.RawData = Encoding.UTF8.GetBytes(postStr);
        //        htpMa.SetHeader("Content-Type", "application/json");
        //        htpMa.Send();


        //}

        public static void SendPost<T>(string url, NetworkHelper.HTTPMethod htpMethod, string postStr, Action<T> YesFun = null, Action NoFun = null) where T : struct
        {
            if (Application.internetReachability == NetworkReachability.NotReachable)
            {
                MetaFramework.UI.UITool.ShowTips("网络错误，请检查网络");
                return;
            }
            NetworkHelper.SendPost(url, htpMethod, postStr, (downloadHandler) =>
            {
                T backData = JsonMapper.ToObject<T>(downloadHandler);
                YesFun?.Invoke(backData);
            }, (code, str) =>
            {
                Debug.Log("code" + code + "str" + str);
                if (code == 429)
                    UITool.ShowTips(str);
                else
                    UITool.ShowTips("服务器维护中，请稍后再试");
                NoFun?.Invoke();
            });
        }


        /// <summary>
        /// 加Herd Token 传输
        /// </summary>
        /// <typeparam name="T"></typeparam>
        /// <param name="url"></param>
        /// <param name="htpMethod"></param>
        /// <param name="headStr"></param>
        /// <param name="postStr"></param>
        /// <param name="YesFun"></param>
        /// <param name="NoFun"></param>
        //public static void SendPost<T>(string url, HTTPMethods htpMethod, string headStr,string postStr, Action<T> YesFun = null, Action NoFun = null) where T : struct
        //{
        //    if (Application.internetReachability == NetworkReachability.NotReachable)
        //    {
        //        MetaFramework.UI.UITool.ShowTips("网络错误，请检查网络");
        //        return;
        //    }
        //        //var htpMa = new HTTPRequest(new Uri(url), htpMethod, ((HTTPRequest httpRequest, HTTPResponse httpResponse) =>
        //        //{
        //        //    if (httpResponse.IsSuccess && httpResponse.StatusCode == 200)
        //        //    {
        //        //        T backData = JsonMapper.ToObject<T>(httpResponse.DataAsText);
        //        //        YesFun?.Invoke(backData);
        //        //    }
        //        //    else
        //        //    {
        //        //        NoFun?.Invoke();
        //        //    }
        //        //}));

        //    //htpMa.SetHeader("Authorization", "Bearer " + headStr);

        //    //if (!string.IsNullOrEmpty(postStr))
        //    //{
        //    //    htpMa.RawData = Encoding.UTF8.GetBytes(postStr);
        //    //    htpMa.SetHeader("Content-Type", "application/json");
        //    //}
        //    //htpMa.Send();

        //}

        public static void SendPost<T>(string url, NetworkHelper.HTTPMethod htpMethod, string headStr, string postStr, Action<T> YesFun = null, Action NoFun = null) where T : struct
        {
            if (Application.internetReachability == NetworkReachability.NotReachable)
            {
                MetaFramework.UI.UITool.ShowTips("网络错误，请检查网络");
                return;
            }
            if (InitJEngine.Instance.debug)
                Debug.Log(url + " + " + postStr);
            NetworkHelper.SendPost(url, htpMethod, headStr, postStr, (downloadHandler) =>
            {
                if (string.IsNullOrEmpty(downloadHandler))
                {
                    NoFun?.Invoke();
                    return;
                }
                T backData = JsonMapper.ToObject<T>(downloadHandler);
                YesFun?.Invoke(backData);
            }, (code, str) =>
            {
                if (code == 429)
                    UITool.ShowTips(str);
                else
                    UITool.ShowTips("服务器维护中，请稍后再试");

                NoFun?.Invoke();
            });
        }


        public static void GetVersion(Action UpCall, Action GoCall)
        {
            if (Application.internetReachability == NetworkReachability.NotReachable)
            {
                MetaFramework.UI.UITool.ShowTips("无网络，请检查网络");
                return;
            }
            string _version = string.Empty;
            string url = $"{Updater.ResourceUrl}/{AssetMgr.GetPlatform}/{Updater.MainPackageName}/PackageManifest_Main.version";
            //string url = "http://121.196.233.214/mowtest/Android/Main/PackageManifest_Main.version";
            var htpMa = new HTTPRequest(new Uri(url), HTTPMethods.Get, ((HTTPRequest httpRequest, HTTPResponse httpResponse) =>
            {
                if (httpResponse.IsSuccess)
                {
                    _version = httpResponse.DataAsText;
                    string nowVersion = PlayerPrefsExtend.GetString("Version");
                    Debug.Log("_version  " + _version);
                    Debug.Log("nowVersion  " + nowVersion);
                    if (!string.IsNullOrEmpty(nowVersion) && !string.IsNullOrEmpty(_version) && _version != nowVersion)
                    {
                        UpCall?.Invoke();
                    }
                    else
                    {
                        GoCall?.Invoke();
                    }
                }
                else
                {
                    MetaFramework.UI.UITool.ShowTips("网络异常，请检查网络");
                }
            }));
            htpMa.Send();
        }

        /// <summary>
        /// 
        /// </summary>
        /// <param name="type">GET或POST</param>
        /// <param name="Url"></param>
        /// <param name="postDataStr"></param>
        /// <returns></returns>
        //public static string SendPost(string type, string Url, string postDataStr,Action<string> Fun = null) //请求方式GET或POST "POST", 
        //{
        //    Debug.Log("post Url = " + Url + " postDataStr = " + postDataStr);
        //    try
        //    {
        //        HttpWebRequest request = (HttpWebRequest)WebRequest.Create(Url);
        //        request.Method = type;
        //        //Post请求方式
        //        request.ContentType = "application/x-www-form-urlencoded";          // 内容类型

        //        //request.
        //        if (!string.IsNullOrEmpty(postDataStr))
        //        {
        //            byte[] payload = Encoding.UTF8.GetBytes(postDataStr);               //将URL编码后的字符串转化为字节
        //            request.ContentLength = payload.Length;                             //设置请求的 ContentLength 
        //            Stream writer = request.GetRequestStream();                         //获得请 求流
        //            writer.Write(payload, 0, payload.Length);                           //将请求参数写入流
        //            writer.Close();                                                     // 关闭请求流
        //        }
        //        HttpWebResponse response = (HttpWebResponse)request.GetResponse();  // 获得响应流
        //        Stream myResponseStream = response.GetResponseStream();
        //        StreamReader myreader = new StreamReader(response.GetResponseStream(), Encoding.UTF8);
        //        string retString = myreader.ReadToEnd();
        //        myreader.Close();
        //        myResponseStream.Close();

        //        if (Fun != null)
        //        {
        //            Fun(retString);
        //        }
        //        return retString;
        //    }
        //    catch (Exception e)
        //    {
        //        Log.PrintError(e);
        //        return null;
        //    }
        //}

        //public void CreateImageRequest( string Url, OnRequestFinishedDelegate CallBack) //type 1:Get, 2:Post
        //{
        //    RequestData data = new RequestData();
        //    var request = HTTPManager.SendRequest(Url, CallBack);
        //    request.Tag = data;
        //    request.Send();
        //    m_requestSave[m_nextGUID] = request;
        //}

        /// <summary>
        /// 
        /// </summary>
        /// <param name="type"> 1:Get, 2:Post</param>
        /// <param name="Url"></param>
        /// <returns></returns>
        //public long CreateUrlRequest(int type, string Url) //type 1:Get, 2:Post
        //{
        //    RequestData data = new RequestData();
        //    data.guid = ++ m_nextGUID;
        //    HTTPMethods t = HTTPMethods.Get;
        //    if (type == 2)
        //        t = HTTPMethods.Post;
        //    var request = new HTTPRequest(new Uri(Url), t, OnDownLoadUrl);
        //    request.Tag = data;
        //    request.Send();

        //    m_requestSave[m_nextGUID] = request;

        //    return m_nextGUID;
        //}

        //void OnDownLoadUrl(HTTPRequest request, HTTPResponse resp)
        //{
        //    var requestdata = (RequestData)request.Tag;

        //    m_requestSave.Remove(requestdata.guid);

        //    //需要检测错误状态
        //    if (resp != null)
        //    {
        //        if (resp.IsSuccess)
        //        {
        //            if (OnDownLoadUrlFinish != null)
        //            {
        //                OnDownLoadUrlFinish(requestdata.guid, resp.DataAsText);
        //            }
        //        }
        //        else
        //        {
        //            if (OnDownloadUrlError != null)
        //                OnDownloadUrlError(requestdata.guid, resp.Message, resp.StatusCode);
        //        }
        //    }
        //    else if (request.Exception != null)
        //    {
        //        int errorCode = 0;
        //        if (request.Exception is SocketException)
        //        {
        //            SocketException socketError = request.Exception as SocketException;
        //            errorCode = socketError.ErrorCode;
        //        }
        //        if (OnDownloadUrlError != null)
        //            OnDownloadUrlError(requestdata.guid, request.Exception.Message, errorCode);
        //    }
        //}
    }
}

