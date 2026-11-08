using HotUpdateScripts;
using HotUpdateScripts.ProtoM;

using JEngine.Core;
using JEngine.Helper;

using MetaFramework.UI;

using System;
using System.Collections.Generic;

using UnityEngine;
using UnityEngine.UI;

using UpdateMode = JEngine.Core.Updater.UpdateMode;

public class UIModelLogin : IUIModel
{

}

public class UIViewLogin : MonoBehaviour, IUIView
{
    public Button Btn_GongGao;
    public Button btn_ChatCustomer;


    public Button Btn_Login_Accout;
    public Button Btn_Login_Phone;
    public InputField intput_account;
    public Text txt_test;
    public Text txt_Version;

    public Button btn_WXLogin;

    public Toggle toggle_Arg;
    public Button btn_Argee;
    public Button btn_YinSi;

    public Image img_SeverNew;
    public Button btn_SeverChoice;
    public Image img_Sever;
    public Text txt_severName;
    public Text txt_SeverDes;
    public Text txt_NewSever;

    public Button Btn_GetSever;

    public GameObject obj_Root;

    public Text txt_Load;
    public Slider slider_Load;

    public void Awake()
    {
    }

    public void Init(GameObject go)
    {
    }
    private void OnDestroy()
    {
        UITool.RemoveUI<UILogin>();
    }
}

public class UILogin : UIBase<UIModelLogin, UIViewLogin>
{
    public override UILayer GetLayer()
    {
        return UILayer.Base;
    }

    public struct Event_Test { }
    public override void OnInit()
    {
        GetView().Btn_GongGao.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UINotice>();
        });

        GetView().btn_ChatCustomer.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            SDKCtr.Instance.OnWXCustomerService("customer service", null, null);
        });

        GetView().Btn_Login_Accout.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            CheakUpData(() =>
            {
                SDKCtr.Instance.OnUpCheak(OnHttp_Accout_Login);
            });

            Dictionary<string, object> properties = new Dictionary<string, object>();
            properties["page"] = "登录界面";
            properties["component"] = "";
            properties["position"] = "";
            properties["event_click"] = "开始游戏";
            TDAnalyticsUtil.Track("click_event", properties, true);

        });
        GetView().Btn_Login_Phone.onClick.AddListener(PhoneLogin);

        GetView().btn_SeverChoice.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIChoiceSever>(null, () => { OnSeverInfoShow(false); });

            Dictionary<string, object> properties = new Dictionary<string, object>();
            properties["page"] = "登录界面";
            properties["component"] = "";
            properties["position"] = "";
            properties["event_click"] = "切服";
            TDAnalyticsUtil.Track("click_event", properties, true);
        });

        AudioManager.Instance.PlayBgMusic(AudioKey.Sound.BG_Login);
        RedPointManager.Instance.ClearRed();

        GetView().txt_Version.text = PlayerPrefsExtend.GetString("Version");

        GetView().btn_WXLogin.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            if (PlayerPrefsExtend.GetInt("Agree") == 0)
            {
                UITool.ShowTips("请先阅读并同意《用户协议》、《隐私政策》");
                return;
            }
        });

        GetView().btn_Argee.onClick.AddListener(delegate
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UILoginDes>(null, UILoginDes.DesType.Agree);
        });

        GetView().btn_YinSi.onClick.AddListener(delegate
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UILoginDes>(null, UILoginDes.DesType.Privacy);
        });

        GetView().toggle_Arg.onValueChanged.AddListener((isOn) =>
        {
            PlayerPrefsExtend.SetInt("Agree", isOn ? 1 : 0);
        });

        GetView().Btn_GetSever.onClick.AddListener(() =>
        {
            if (string.IsNullOrEmpty(GetView().intput_account.text))
            {
                UITool.ShowTips("请输入账号");
                return;
            }
            PlayerPrefsExtend.SetString("TextAcount", GetView().intput_account.text);
            GetSeverInfo();
        });

        if (GameM.Instance.BaseScale == 0)
        {
            GameM.Instance.BaseScale = Mathf.Max((uiGo.GetComponent<RectTransform>().rect.width / 720f), (uiGo.GetComponent<RectTransform>().rect.height / 1280f));
        }
        OnInitUserData();
    }

    public override void OnEnter(params object[] args)
    {
        GetView().intput_account.text = PlayerPrefsExtend.GetString("TextAcount");

        OnSeverBaseGet();

    }

    /// <summary>
    /// 获取服务器信息
    /// </summary>
    private void OnSeverBaseGet()
    {
        GetView().intput_account.gameObject.SetActive(InitJEngine.Instance.IsFreeLogin);

        GetView().obj_Root.SetActive(ExcelDataManager.Instance.IsLoad);
        GetView().slider_Load.gameObject.SetActive(false);

        if (SDKCtr.Instance.CurrentPlatform == SDKCtr.PlatformType.WebGL_WX)
        {
            GetView().obj_Root.SetActive(ExcelDataManager.Instance.IsLoad);
            GetView().slider_Load.gameObject.SetActive(false);

            UITool.Open<UILoading>((ui) =>
            {
                //微信code获取
                SDKCtr.Instance.OnWXWebLogin((code) =>
                {
                    GameM.Instance.WXCode = code;
                    UITool.Open<UILoading>((ui) =>
                    {
                        string clickid = GameUtil.Instance.GetWXLaunchParamClickId();
                        //换取 OpenID
                        HttpC.OnGetOpenID(code, clickid, (msg) =>
                        {
                            GameM.Instance.WXOpneID = msg.data.openid;
                            PlayerPrefsExtend.SetString("TextAcount", GameM.Instance.WXOpneID);
                            GetView().intput_account.text = GameM.Instance.WXOpneID;
                            SDKCtr.Instance.DNSetOpenId(msg.data.openid);
                            if (PlayerPrefsExtend.GetInt(GameConst.IsFirstOpenMiniGame) == 0)
                            {
                                SDKCtr.Instance.DNOnRegister();
                                PlayerPrefsExtend.SetInt(GameConst.IsFirstOpenMiniGame, 1);
                            }
                            GetSeverInfo();
                            TDAnalyticsUtil.Login(GameM.Instance.WXOpneID);
                            TDAnalyticsUtil.SetDistinctId(GameM.Instance.WXOpneID);
                            Debug.Log("TDAnalyticsUtil.Login() ----");
                        });
                    }, new UILoading.Loading { _time = 0, _str = "信息校验..." });

                    SDKCtr.Instance.OnAuthorize(null, null);
                }, null);
            }, new UILoading.Loading { _time = 0, _str = "信息获取..." });
        }
        else if (SDKCtr.Instance.CurrentPlatform == SDKCtr.PlatformType.Android)
        {
            GetView().obj_Root.SetActive(false);
            UITool.Open<UIBQSDKMgr>();
        }
        else
        {
            GetView().obj_Root.SetActive(true);
            GetView().slider_Load.gameObject.SetActive(false);
            GetSeverInfo();
        }
    }

    private void GetSeverInfo()
    {
        // Local single-server: no remote go-admin request.
        // Without this SeverInfoDic stays empty and login throws KeyNotFoundException.
        GameM.Instance.SeverIdx = 0;
        GameM.Instance.SeverDicIdx = 0;
        GameM.Instance.InitLocalSeverInfo();
        OnSeverInfoShow(false);
    }

    /// <summary>
    /// 显示服务器信息
    /// </summary>
    /// <param name="isInit"></param>
    public void OnSeverInfoShow(bool isInit)
    {
        GetView().txt_NewSever.text = GameM.Instance.SeverInfoDic[1][0].name;
        ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CreatRoleAtlas, "img_" + GameM.Instance.SeverInfoDic[1][0].status, GetView().img_SeverNew);

        if (GameM.Instance.SeverInfoDic[0].Count == 0 && isInit)
        {
            GameM.Instance.SeverIdx = 0;
            GameM.Instance.SeverDicIdx = 1;
        }
        var SelectSever = GameM.Instance.SeverInfoDic[GameM.Instance.SeverDicIdx][GameM.Instance.SeverIdx];
        GetView().txt_severName.text = SelectSever.name;

        Debug.Log("SelectSever.ipaddr --- " + SelectSever.ipaddr);

        ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CreatRoleAtlas, "img_" + SelectSever.status, GetView().img_Sever);
    }

    /// <summary>
    /// 登录账号（本地服：跳过远端版本校验，直接登录）
    /// </summary>
    private void OnHttp_Accout_Login()
    {
        if (SDKCtr.Instance.CurrentPlatform == SDKCtr.PlatformType.WebGL_WX && GameM.Instance.SeverInfoDic[GameM.Instance.SeverDicIdx][GameM.Instance.SeverIdx].status == 0)
        {
            UITool.ShowTips("服务器已关闭");
            GetView().obj_Root.SetActive(true);
            return;
        }

        UITool.Open<UILoading>((ui) =>
        {
            var loginRequest = new HTTP_LoginRequest
            {
                appid = "wxChat",
                channelId = PlayerPrefsExtend.GetString(SDKCtr.SaveKey.WX_AdID.ToString()),
                imei = PlayerPrefsExtend.GetString("TextAcount"),
                serverId = GameM.Instance.SeverInfoDic[GameM.Instance.SeverDicIdx][GameM.Instance.SeverIdx].server_id,
                clickid = GameUtil.Instance.GetWXLaunchParamClickId(),
            };
            HttpC.OnLoginSever(loginRequest, HttpC.LoginType.Account, OnLoginUIShow, () => { GetView().obj_Root.SetActive(true); });
        }, new UILoading.Loading { _time = 0, _str = "服务器链接中..." });
    }

    //表格数据加载完成后，显示登录界面
    private void OnLoginUIShow()
    {
        void OnLoginShow()
        {
            void OnShowLogin()
            {
                GetView().obj_Root.SetActive(ExcelDataManager.Instance.IsLoad);
                GetView().slider_Load.gameObject.SetActive(!ExcelDataManager.Instance.IsLoad);
            }
            if (RoleM.Instance.LoginCharacterInfo.Count == 0)
            {
                UITool.Open<UICreatRole>((ui) => { OnShowLogin(); });
            }
            else
            {
                UITool.Open<UISelectRole>((ui) => { OnShowLogin(); });
            }
        }
        if (!ExcelDataManager.Instance.IsLoad)
        {
            GetView().obj_Root.SetActive(ExcelDataManager.Instance.IsLoad);
            GetView().slider_Load.gameObject.SetActive(!ExcelDataManager.Instance.IsLoad);
            ExcelDataManager.Instance.LoadAllCFG((loadCount, totalCount) =>
            {
                float valuepro = loadCount / (float)totalCount;
                GetView().txt_Load.text = $"登录中...{GameUtil.Instance.GetPercentStr(valuepro)}";
                GetView().slider_Load.value = valuepro;
                if (loadCount == totalCount)
                {
                    OnLoginShow();
                }
            });
        }
        else
        {
            OnLoginShow();
        }
    }

    private void PhoneLogin()
    {
        //UITool.Open<UIPhoneLogin>();
    }

    private void OnInitUserData()
    {
        RoleM.Instance.NowOpenLevel = 0;
        RoleM.Instance.NowOpenJob = 0;
        RoleM.Instance.IsReceivedGuideData = false;
        RoleM.Instance.IsFirstReceiveBaseInfo = true;
        StateM.Instance.isSendUpgradeOrSpeedUpBox = false;
        StateM.Instance.isHasRefineResult = false;
        GuidanceM.Instance.OnInit();
        ChatM.Instance.ClearAllChat();
        RoleM.Instance.TotalCombat = 0;
        RoleM.Instance.refineResuleCache = null;
    }

    //更新检测
    private void CheakUpData(Action BackCall)
    {
        if (Updater.Mode != UpdateMode.WebGLRemote)
        {
            Debug.Log("CheakUpData + " + Updater.Mode);
            BackCall?.Invoke();
            return;
        }

        UITool.Open<UILoading>((ui) =>
        {
            AssetMgr.CheakUpdatePackage((versionStr) =>
            {
                Debug.Log("versionStr + " + versionStr);
                UITool.Close<UILoading>();
                if (PlayerPrefsExtend.GetString("Version") != versionStr)
                    OnWXRestart();
                else
                    BackCall?.Invoke();
            }, () =>
            {
                Debug.Log("Err + ");
            });

        }, new UILoading.Loading { _time = 0, _str = "资源校验中" });
    }

    private void OnWXRestart()
    {
        UITool.Open<UIMessageBoxOne>(null, new UIMessageBoxOne.CallBack
        {
            showStr = "新版本已经准备好，是否重启应用？",
            endFun = delegate
            {
                SDKCtr.Instance.OnWXRestartMiniProgram((str) =>
                {
                    UITool.Open<UIMessageBoxOne>(null, new UIMessageBoxOne.CallBack
                    {
                        showStr = "重启失败，请手动重启",
                    });
                });
            },
        });
    }



}

