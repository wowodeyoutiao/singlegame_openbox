using JEngine.LitJson;

using System;
using System.Collections;
using System.Collections.Generic;
using System.Runtime.InteropServices;

using ThinkingData.Analytics;

using UnityEngine;
using UnityEngine.Networking;

using WeChatWASM;

public class SDKCtr : MonoBehaviour
{
    public static SDKCtr Instance
    {
        get; set;
    }

    [Header("默认交互 包.静态类名(AndroidJavaClass)")]
    [SerializeField] string defaultAndroidPackageClass = "com.zsjzz.east.AndroidUnityConnecter";
    [Header("默认 静态类.接收消息方法名")]
    [SerializeField] string defaultReceiveMethodName = "ReceiveUnityMsg";

    private AndroidJavaClass mainActivityClass = null;
    private AndroidJavaObject mainActivityObject = null;

    public readonly string WXAPPID = "wx9d99386ab1e0721d";
    private readonly string WXSECRET = "a9e59e969eedeb31e0ceb2c0800862ad";

    public enum PlatformType
    {
        Android,
        IOS,
        WebGL,
        WebGL_WX,
        Editor
    }

    /// <summary>
    /// 当前运行平台类型
    /// </summary>
    public PlatformType CurrentPlatform
    {
        get
        {
#if UNITY_ANDROID
            return PlatformType.Android;
#elif UNITY_IOS
            return PlatformType.IOS;
#elif !UNITY_EDITOR && UNITY_WEBGL && WEIXINMINIGAME
            return PlatformType.WebGL_WX;
#elif !UNITY_EDITOR && UNITY_WEBGL
            return PlatformType.WebGL;
#endif
            return PlatformType.Editor;
        }
    }
    public enum SaveKey
    {
        WX_OPEN_ID,
        WX_ACCESS_TOKEN,
        WX_SESSION_KEY,
        WX_REFRESH_TOKEN,
        WX_REFRESH_TOKEN_TIME,

        WX_Platform,
        WX_AdID,
    }

    private void Awake()
    {
        Debug.Log("SDK Awake");
        if (Instance != null)
            return;
        Instance = this;
        DontDestroyOnLoad(this);

#if !UNITY_EDITOR && UNITY_WEBGL && WEIXINMINIGAME
Debug.Log(" WebGLInput.mobileKeyboardSupport");
        WebGLInput.mobileKeyboardSupport = true;

        InitWXWebSDK();

#endif

    }


    void Start()
    {
        Debug.Log("SDK Start");

#if UNITY_ANDROID
        mainActivityClass = new AndroidJavaClass(defaultAndroidPackageClass);
        mainActivityObject = new AndroidJavaObject(defaultAndroidPackageClass);

#elif WEIXINMINIGAME && UNITY_WEBGL && !UNITY_EDITOR
        // 微信小游戏环境下不需要初始化AndroidJavaClass
        Debug.Log("SDK WEIXINMINIGAME");
#endif

    }

    /// <summary>
    /// 设置屏幕常亮选项
    /// </summary>
    private void OnKeepScreenOnOption()
    {
        SetKeepScreenOnOption option = new SetKeepScreenOnOption
        {
            keepScreenOn = true,
            success = (result) => {
                Debug.Log("SetKeepScreenOnOption success" + result.errMsg);

                WX.OnHide((result) =>
                {
                    Debug.Log("WX.OnHide");
                    WXOnHide();
                });
            },
            fail = (result) => { Debug.LogError("SetKeepScreenOnOption " + result.errMsg); }
        };
        WX.SetKeepScreenOn(option);
    }

    void OnEnable()
    {
#if WEIXINMINIGAME && UNITY_WEBGL
        WX.OnShow((result) =>
        {
            Debug.Log("WX.OnShow  SetKeepScreenOnOption");
            OnKeepScreenOnOption();
            OnMiniGameShow();
            if (result.query.ContainsKey("weixinadinfo"))
            {
                string aid = result.query["weixinadinfo"].Split('.')[0];
                PlayerPrefsExtend.SetString(SaveKey.WX_AdID.ToString(), aid);
            }
        });
#else
        Screen.sleepTimeout = SleepTimeout.NeverSleep;
#endif
    }

    void OnDisable()
    {
#if WEIXINMINIGAME && UNITY_WEBGL && !UNITY_EDITOR
        WX.OnHide((result) =>
        {
            Debug.Log("WX.OnHide");
            WXOnHide();
        });
#endif

    }

    #region WXWeb
    Action<string> LoginCallBack;
    private void InitWXWebSDK()
    {
        WX.InitSDK((code) =>
        {

            Debug.Log("微信SDK初始化成功" + _AccountInfo.miniProgram.envVersion);
            OnKeepScreenOnOption();
            InitRewardedVideo();
            SDKCtr.Instance.OnUpCheak(null);
            OnWXWebSystemInfo(null);
            ShuShuInit();
            //OnMiniGameShare();
            WxLaunch();
            OnUpCheak(null);
            SendTDDeviceActivation();
            TDAnalytics.Flush();
        });
    }

    private void OnMiniGameShow()
    {
        entersForegroundTime = DateTime.Now;
        bool isFirst = PlayerPrefsExtend.GetInt("shushu_ta_mp_show") == 0;
        TDAnalytics.Track("ta_mp_show", new Dictionary<string, object>
        {
            ["is_first_time"] = isFirst
        });

        if (PlayerPrefsExtend.GetInt("shushu_ta_mp_show") == 0)
        {
            PlayerPrefsExtend.SetInt("shushu_ta_mp_show", 1);
        }
    }


    private void OnMiniGameShare()
    {
        WX.OnShareTimeline(WXOnShareTimeline);
    }

    /// <summary>
    /// 微信Web登录
    /// </summary>
    public void OnWXWebLogin(Action<string> _wxCallBack, Action _FailCall)
    {
        LoginCallBack = _wxCallBack;

        LoginOption loginOption = new LoginOption()
        {
            success = LoginSuccess,
            fail = (result) => { Debug.Log($"登录失败 code={result.errno} msg={result.errMsg}"); _FailCall?.Invoke(); }
        };
        Debug.Log("开始登录");
        WX.Login(loginOption);
    }

    /// <summary>
    /// 登陆成功回调
    /// </summary>
    /// <param name="result"></param>
    private void LoginSuccess(LoginSuccessCallbackResult result)
    {
        Debug.Log($"登录成功 code={result.code} msg={result.errMsg}");
        LoginCallBack?.Invoke(result.code);
    }

    /// <summary>
    /// 微信Web支付（虚拟支付）
    /// </summary>
    /// <param name="paySig"></param>
    /// <param name="signData"></param>
    /// <param name="WinAc"></param>
    /// <param name="FailAc"></param>
    public void OnWXWebPay(string paySig, string signData, string signature, Action WinAc, Action FailAc)
    {
        Debug.Log("微信Web支付 " + paySig + "   " + signData + "   " + signature);
        RequestMidasPaymentGameItemOption requestMidasPaymentGameItemOption = new RequestMidasPaymentGameItemOption()
        {
            signData = signData, // 这里需要替换为实际的签名数据
            paySig = paySig, // 这里需要替换为实际的支付签名
            signature = signature, // 这里需要替换为实际的签名
            success = (result) => { Debug.Log($"支付成功 code={result.errCode} msg={result.errMsg}"); WinAc?.Invoke(); },
            fail = (result) => { Debug.Log($"支付失败 code={result.errCode} msg={result.errMsg}"); FailAc?.Invoke(); }
        };
        WX.RequestMidasPaymentGameItem(requestMidasPaymentGameItemOption);
    }

    /// <summary>
    /// 微信客服
    /// </summary>
    /// <param name="sessionFrom"></param>
    public void OnWXCustomerService(string sessionFrom, Action WinAc, Action FailAc)
    {
        if (CurrentPlatform != PlatformType.WebGL_WX)
            return;

        Debug.Log("用微信支付 回话" + sessionFrom);
        WX.OpenCustomerServiceConversation(new OpenCustomerServiceConversationOption
        {
            sessionFrom = sessionFrom,
            showMessageCard = false,
            success = (result) => { Debug.Log("客服会话打开成功"); WinAc?.Invoke(); },
            fail = (result) => { Debug.LogError("客服会话打开失败: " + result.errMsg); FailAc?.Invoke(); }
        });
    }

    public void OnWXWebSystemInfo(Action CallBack)
    {
        Debug.Log("获取微信Web系统信息");
        WX.GetSystemInfo(new GetSystemInfoOption
        {
            success = (result) => {
                Debug.Log($" platform ={result.platform} ");
                PlayerPrefsExtend.SetString(SaveKey.WX_Platform.ToString(),result.platform);
                CallBack?.Invoke();
            },
            fail = (result) => { Debug.Log($"信息获取失败 " + result.errMsg); }
        });
    }


    private AccountInfo accountInfo = null;
    private AccountInfo _AccountInfo
    {
        get
        {
            if (accountInfo == null)
            {
                accountInfo = WX.GetAccountInfoSync();
            }
            return accountInfo;
        }
    }
    public void OnAuthorize(Action WinCall, Action FailCall)
    {
        if (CurrentPlatform != PlatformType.WebGL_WX)
        {
            WinCall?.Invoke();
            return;
        }
        GetPrivacySettingOption privacyOption = new GetPrivacySettingOption();
        privacyOption.success = (res) => {
            if (res.needAuthorization)  //隐私协议
            {
                WX.RequirePrivacyAuthorize(new RequirePrivacyAuthorizeOption
                {
                    success = (_) =>
                    {
                        WinCall?.Invoke();
                    },
                    fail = (error) =>
                    {
                        Debug.LogError($"隐私协议授权失败: {error.errMsg}");
                        FailCall?.Invoke();
                    }
                });
            }
            else
            {
                WinCall?.Invoke();
            }
        };
        WX.GetPrivacySetting(privacyOption);
    }

    /// <summary>
    /// 检测升级
    /// </summary>
    /// <param name="NoUpCall"></param>
    public void OnUpCheak(Action NoUpCall)
    {
        if (CurrentPlatform != PlatformType.WebGL_WX)
        {
            NoUpCall?.Invoke();
            return;
        }

        if (!WX.CanIUse("GetUpdateManager"))
        {
            WX.ShowModal(new ShowModalOption
            {
                showCancel = false,
                title = "提示",
                content = "当前微信版本过低，无法使用该功能，请升级到最新微信版本后重试。",
            });
        }
        else
        {
            var updateManager = WX.GetUpdateManager();
            updateManager.OnCheckForUpdate((result) =>
            {
                Debug.Log($"检查更新结果：{result.hasUpdate}");
                if (result.hasUpdate)
                {
                    updateManager.OnUpdateReady((res) => {
                        WX.ShowModal(new ShowModalOption
                        {
                            title = "更新提示",
                            content = "新版本已经准备好，是否重启应用？",
                            showCancel = false,
                            confirmText = "确定",
                            success = (modalResult) =>
                            {
                                updateManager.ApplyUpdate();
                            },
                            fail = (error) => { Debug.LogError($"模态框显示失败: {error.errMsg}"); }
                        });
                    });

                    updateManager.OnUpdateFailed((res) => {
                        WX.ShowModal(new ShowModalOption
                        {
                            title = "已有新版本",
                            content = "请您删除小程序，重新搜索进入",
                            showCancel = false,
                        });
                    });
                }
                else
                {
                    NoUpCall?.Invoke();
                    Debug.Log("当前已是最新版本");
                }
            });
        }
    }

    private WXGameClubButton wXGameClubButton = null;

    /// <summary>
    /// 打开微信游戏圈按钮
    /// </summary>
    /// <param name="topLeftPos"></param>
    /// <param name="gameClubButtonRectTransform"></param>
    /// <param name="isShow"></param>
    public void OpenWXGamePag(Vector2 topLeftPos, RectTransform gameClubButtonRectTransform, bool isShow)
    {
        if (CurrentPlatform != PlatformType.WebGL_WX)
            return;

        if (wXGameClubButton == null && gameClubButtonRectTransform != null)
        {
            Debug.Log("创建游戏圈按钮");

            WeChatWASM.WXCreateGameClubButtonParam param = new WeChatWASM.WXCreateGameClubButtonParam();
            param.type = WeChatWASM.GameClubButtonType.text;
            param.text = "  "; // 
            WeChatWASM.GameClubButtonStyle gameClubButtonStyle = new WeChatWASM.GameClubButtonStyle();

            //设置位置和大小
            //Vector2 topLeftPos = UtilManager.Instance.GetUITopLeftCorner(gameClubButtonRectTransform);
            float hRadion = (float)WeChatWASM.WX.GetSystemInfoSync().screenHeight / Screen.height;
            float wRadion = (float)WeChatWASM.WX.GetSystemInfoSync().screenWidth / Screen.width;

            gameClubButtonStyle.left = (int)(topLeftPos.x * wRadion);
            gameClubButtonStyle.top = (int)(WeChatWASM.WX.GetSystemInfoSync().screenHeight - topLeftPos.y * hRadion);
            gameClubButtonStyle.width = (int)(gameClubButtonRectTransform.sizeDelta.x * wRadion);
            gameClubButtonStyle.height = (int)(gameClubButtonRectTransform.sizeDelta.y * hRadion);

            Debug.Log("gameClubButtonStyle.left " + gameClubButtonStyle.left + "topLeftPos.top" + gameClubButtonStyle.top + "gameClubButtonStyle ." + gameClubButtonStyle.width + "gameClubButtonStyle " + gameClubButtonStyle.height);
            param.style = gameClubButtonStyle;

            //创建游戏圈按钮
            wXGameClubButton = WX.CreateGameClubButton(param);
        }
        if (wXGameClubButton != null)
        {
            if (isShow)
                wXGameClubButton.Show();
            else
                wXGameClubButton.Hide();
        }



        //wXGameClubButton.Hide(); // 隐藏按钮
        //WX.CreateGameClubButton(new WXCreateGameClubButtonParam { 
        //    type = GameClubButtonType.image,
        //    icon = GameClubButtonIcon.green,
        //    style = new GameClubButtonStyle
        //    {
        //        left = 10,
        //        top = 10,
        //        width = 50,
        //        height = 50,
        //    },
        //});

    }

    private WeChatWASM.WXRewardedVideoAd wxRewardedVideoAd;
    /// <summary>
    /// 激励广告初始化
    /// </summary>
    private void InitRewardedVideo()
    {
        if (wxRewardedVideoAd != null) return;
        var param = new WeChatWASM.WXCreateRewardedVideoAdParam
        {
            adUnitId = "adunit-40905fa00039b762",
        };
        wxRewardedVideoAd = WeChatWASM.WXBase.CreateRewardedVideoAd(param);
        wxRewardedVideoAd.OnClose(OnRewardedVideoClose);
    }

    private Action ADWinCall;
    private Action ADFailCall;

    /// <summary>
    /// 激励视频广告展示
    /// </summary>
    /// <param name="winCall"></param>
    /// <param name="fialCall"></param>
    public void ShowRewardedVideo(Action winCall, Action fialCall)
    {
        ADWinCall = winCall;
        ADFailCall = fialCall;
        if (wxRewardedVideoAd != null)
            wxRewardedVideoAd.Show(null, (message) => {
                wxRewardedVideoAd.Load((load) => {
                    wxRewardedVideoAd.Show();
                });
                Debug.Log("观看广告 err  " + message.errMsg);
            });
        else
            fialCall?.Invoke();
    }

    /// <summary>
    /// 激励视频广告关闭回调
    /// </summary>
    /// <param name="obj"></param>
    private void OnRewardedVideoClose(WXRewardedVideoAdOnCloseResponse obj)
    {
        if (obj != null && obj.isEnded) //正常播放结束，可以发放奖励
        {
            ADWinCall?.Invoke();
            // 正常播放结束，可以下发游戏奖励
        }
        else//中途播放结束
        {
            Debug.Log("关闭广告，不发放奖励 ");
            ADFailCall?.Invoke();
        }
    }

    /// <summary>
    /// 重启小游戏
    /// </summary>
    /// <param name="fail_action"></param>
    public void OnWXRestartMiniProgram(Action<string> fail_action)
    {
        WX.RestartMiniProgram(new RestartMiniProgramOption {
            fail = (result) => { Debug.LogError("重启小游戏失败: " + result.errMsg); fail_action?.Invoke(result.errMsg); },
        });
    }

    #endregion


    public void QuitGame()
    {
        WeChatWASM.WX.ExitMiniProgram(new ExitMiniProgramOption
        {
            success = (result) => { Debug.Log("退出小游戏成功"); },
            fail = (result) => {
                WX.ShowModal(new ShowModalOption
                {
                    title = "提示",
                    content = "游戏数据异常",
                    showCancel = false,

                });
            }
        });
    }


    /// <summary>
    /// 文本复制
    /// </summary>
    /// <param name="str"></param>
    /// <param name="WinCall"></param>
    public void CopyString(string str, Action WinCall)
    {
#if WEIXINMINIGAME && UNITY_WEBGL && !UNITY_EDITOR
        OnAuthorize(()=>{
             WX.SetClipboardData(new SetClipboardDataOption {
                data = str,
                success = (result) => { WinCall?.Invoke(); },
                fail = (result) => { Debug.Log("复制失败 "+ result.errMsg); }
            });
        },null);
#else 
        GUIUtility.systemCopyBuffer = str;
        WinCall?.Invoke();
#endif
    }

    /// <summary>
    /// 设置帧率
    /// </summary>
    /// <param name="FpsNum"></param>
    public void SetPreferredFramesPerSecond(int FpsNum)
    {
#if WEIXINMINIGAME && UNITY_WEBGL && !UNITY_EDITOR
        WX.SetPreferredFramesPerSecond(FpsNum);
#else
        Application.targetFrameRate = FpsNum;
#endif
    }

    private Action WXCallBack;
    public void WeChatLogin(Action _wxCallBack)
    {
        Debug.Log("WeChatLogin  WXLogin");
        WXCallBack = _wxCallBack;

        string WX_REFRESH_TOKEN_TIME = PlayerPrefs.GetString(SaveKey.WX_REFRESH_TOKEN_TIME.ToString());
        if (string.IsNullOrEmpty(WX_REFRESH_TOKEN_TIME) || (DateTime.Now - DateTime.Parse(WX_REFRESH_TOKEN_TIME)).TotalDays > 29)
        {
            mainActivityObject.CallStatic("WXLogin", "微信接收到了Login请求");
        }
        else
        {
            GetTokenFromServer();
        }
    }

    public void LoginCallback(string str)
    {
        Debug.Log("登陆返回：" + str);
        GetTokenFromServer(str);
    }

    private void GetTokenFromServer(string code)
    {
        // Debug.Log("code 登陆返回：" + code);
        string url = $"https://api.weixin.qq.com/sns/oauth2/access_token?appid={WXAPPID}&secret={WXSECRET}&code={code}&grant_type=authorization_code";
        StartCoroutine(GetTokenByWWW(url, false));
    }

    private void GetTokenFromServer()
    {
        //Debug.Log("刷新Token");
        string url = $"https://api.weixin.qq.com/sns/oauth2/refresh_token?appid={WXAPPID}&grant_type=refresh_token&refresh_token={PlayerPrefs.GetString(SaveKey.WX_REFRESH_TOKEN.ToString())}";
        StartCoroutine(GetTokenByWWW(url, false));
    }

    [System.Serializable]
    private class WXLoginInfo
    {
        public string access_token;
        public int expires_in;
        public string refresh_token;
        public string openid;
        public string scope;
        public string unionid;
        public int errcode;
        public string errmsg;
        public string session_key;
    }

    private IEnumerator GetTokenByWWW(string url, bool isRefresh)
    {
        UnityWebRequest webRequest = UnityWebRequest.Get(url);
        webRequest.downloadHandler = (DownloadHandler)new DownloadHandlerBuffer();
        yield return webRequest.SendWebRequest();
        if (webRequest.result == UnityWebRequest.Result.ConnectionError || webRequest.result == UnityWebRequest.Result.ProtocolError)
        {
            SendMsgToAndroid("toast", "网络错误");
            yield break;
        }

        yield return webRequest.isDone;
        var info = JsonMapper.ToObject<WXLoginInfo>(webRequest.downloadHandler.text);
        if (!string.IsNullOrEmpty(info.errmsg))
        {
            SendMsgToAndroid("toast", "get token error " + info.errcode);
            yield break;
        }
        PlayerPrefs.SetString(SaveKey.WX_ACCESS_TOKEN.ToString(), info.access_token);
        PlayerPrefs.SetString(SaveKey.WX_SESSION_KEY.ToString(), info.session_key);
        PlayerPrefs.SetString(SaveKey.WX_OPEN_ID.ToString(), info.openid);
        PlayerPrefs.SetString(SaveKey.WX_REFRESH_TOKEN.ToString(), info.refresh_token);
        PlayerPrefs.SetString(SaveKey.WX_REFRESH_TOKEN_TIME.ToString(), DateTime.Now.ToString());

        WXCallBack?.Invoke();
    }

    public void ReceiveAndroidMsg(string args)
    {
        //根据接收到的消息自定义处理
        Debug.Log($"来自Android的消息:{args}");
    }

    public void TextMsg(string textMsg)
    {
        Debug.Log($"来自textMsg的消息:{textMsg}");
    }

    /// <summary>
    /// 发送消息到Android
    /// </summary>
    /// <param name="key">key</param>
    /// <param name="args">参数</param>
    public void SendMsgToAndroid(string key, string args)
    {
#if UNITY_ANDROID
        mainActivityObject.CallStatic(defaultReceiveMethodName, key, args);
#endif
    }

    private Action ALiPayWinBack;
    private Action ALiPayFailBack;
    //支付调用
    public void OnALiPay(string order, Action _winCall, Action _failCall)
    {
        Debug.Log("支付调起 ");
        ALiPayWinBack = _winCall;
        ALiPayFailBack = _failCall;

        if (!IsZhiFuBaoAvilible())
        {
            _failCall?.Invoke();
        }
        else
        {
            mainActivityObject.CallStatic("ALiPay", order);
        }
    }

    //支付返回
    public void ALiPayResultCallBack(string result)
    {
        Debug.Log("支付回调信息  " + result);

        string memo = ExtractValue(result, "memo=");
        if (!string.IsNullOrEmpty(memo))
            SendMsgToAndroid("toast", memo);

        if (ExtractValue(result, "resultStatus=") == "9000")
        {
            ALiPayWinBack?.Invoke();
        }
        else
        {
            ALiPayFailBack?.Invoke();
        }
    }

    private string ExtractValue(string input, string key)
    {
        // 查找键的位置
        int keyIndex = input.IndexOf(key);
        if (keyIndex == -1)
        {
            return ""; // 键不存在
        }

        // 找到键后面的花括号
        int startIndex = keyIndex + key.Length;
        int endIndex = input.IndexOf('}', startIndex);
        if (endIndex == -1)
        {
            return ""; // 没有找到结束的花括号
        }

        // 提取值
        string value = input.Substring(startIndex, endIndex - startIndex).Trim('{', '}');
        return value;
    }


    /// <summary>
    /// 支付宝
    /// </summary>
    /// <returns></returns>
    public bool IsZhiFuBaoAvilible()
    {
#if UNITY_STANDALONE_WIN || UNITY_STANDALONE_OSX || UNITY_EDITOR
        return false;
#elif UNITY_ANDROID
        return mainActivityObject.CallStatic<bool>("isZhiFuBaoAvilible");
#elif UNITY_IOS
        return false;
#elif UNITY_WEBGL
        return false;
#endif
    }

    #region shushuSDK

    // 本次使用时长，单位秒
    long durationRun;
    DateTime entersForegroundTime;
    [Space(10)]
    public string shushuAppId = "36b863e73f2347e5810be7ea42137328";
    public string shushuTestAppId = "debug-appid";
    public string shushuUrl = "https://shushu.shbinqigame.com/";
    public bool shushuTestMode = false;

    private void ShuShuInit()
    {
        string appId = shushuTestMode ? shushuTestAppId : shushuAppId;
        TDConfig tdConfig = new TDConfig(appId, shushuUrl);
        TDAnalytics.Init(tdConfig);
        Debug.Log("-- DAnalytics.Init()");
    }

    private void WxLaunch()
    {
        bool isFirst = PlayerPrefsExtend.GetInt("shushu_ta_mp_launch") == 0;
        TDAnalytics.Track("ta_mp_launch", new Dictionary<string, object>
        {
            ["is_first_time"] = isFirst
        });

        if (PlayerPrefsExtend.GetInt("shushu_ta_mp_launch") == 0)
        {
            PlayerPrefsExtend.SetInt("shushu_ta_mp_launch", 1);
        }
    }

    private void WXOnHide()
    {
        durationRun = (long)Math.Floor((DateTime.Now - entersForegroundTime).TotalSeconds);
        TDAnalytics.Track("ta_mp_hide", new Dictionary<string, object>
        {
            ["#duration"] = durationRun
        });
    }

    private void WXOnShareTimeline(Action<OnShareTimelineListenerResult> action)
    {
        bool isFirst = PlayerPrefsExtend.GetInt("shushu_ta_mp_show") == 0;
        TDAnalytics.Track("ta_mp_share", new Dictionary<string, object>
        {
            ["is_first_time"] = isFirst
        });

        if (PlayerPrefsExtend.GetInt("shushu_ta_mp_share") == 0)
        {
            PlayerPrefsExtend.SetInt("shushu_ta_mp_share", 1);
        }
    }

    /// <summary>
    /// 设备首次启动时，使用首次事件上报
    /// </summary>
    public void SendTDDeviceActivation()
    {
        TDFirstEventModel firstEvent = new TDFirstEventModel("new_device");
        TDAnalytics.Track(firstEvent);
    }

    public Dictionary<string, string> GetWXLaunchParams()
    {
#if !UNITY_EDITOR && UNITY_WEBGL && WEIXINMINIGAME
        return WX.GetLaunchOptionsSync().query;
#else
        return new Dictionary<string, string>();
#endif
    }

    public Dictionary<string, string> GetWXEnterParams()
    {
#if !UNITY_EDITOR && UNITY_WEBGL && WEIXINMINIGAME
        return WX.GetEnterOptionsSync().query;
#else
        return new Dictionary<string, string>();
#endif
    }

    #endregion

    #region 腾讯广告上报
    public void DNSetOpenId(string appId)
    {
#if !UNITY_EDITOR && UNITY_WEBGL && WEIXINMINIGAME
        SetOpenId(appId);
#endif
    }

    public void DNSetUnionId(string unionid)
    {
#if !UNITY_EDITOR && UNITY_WEBGL && WEIXINMINIGAME
        SetUnionId(unionid);
#endif
    }

    public void DNOnRegister()
    {
#if !UNITY_EDITOR && UNITY_WEBGL && WEIXINMINIGAME
        OnRegister();
#endif
    }

    public void DNCreateRole(string userName)
    {
#if !UNITY_EDITOR && UNITY_WEBGL && WEIXINMINIGAME
        CreateRole(userName);
#endif
    }

    /// <summary>
    /// 支付上报
    /// </summary>
    /// <param name="price">价格:分</param>
    public void DNPurchase(int price)
    {
#if !UNITY_EDITOR && UNITY_WEBGL && WEIXINMINIGAME
        Purchase(price);
#endif
    }

    /// <summary>
    /// 引导完成
    /// </summary>
    public void DNTutorialFinish()
    {
#if !UNITY_EDITOR && UNITY_WEBGL && WEIXINMINIGAME
        TutorialFinish();
#endif
    }

#if !UNITY_EDITOR && UNITY_WEBGL && WEIXINMINIGAME

    [DllImport("__Internal")]
    private static extern void SetOpenId(string appId);

    [DllImport("__Internal")]
    private static extern void SetUnionId(string unionid);

    [DllImport("__Internal")]
    private static extern void OnRegister();

    [DllImport("__Internal")]
    private static extern void CreateRole(string userName);

    [DllImport("__Internal")]
    private static extern void Purchase(int price);

    [DllImport("__Internal")]
    private static extern void TutorialFinish();
#endif

    #endregion
}
