using Cysharp.Threading.Tasks;
using HotUpdateScripts;
using HotUpdateScripts.ProtoM;
using JEngine.Core;
using JEngine.Event;
using JEngine.Helper;
using JEngine.MyNet;
using MetaFramework.UI;
using pb;
using System;
using System.Collections.Generic;
//using System.Threading.Tasks;
using UnityEngine;
using UnityEngine.SceneManagement;

public class TimeItem
{
    public TimeItem()
    {
    }

    public TimeItem(float _DelayTime, Action _action) : this()
    {
        DelayTime = _DelayTime;
        action = _action;
        IsUsed = false;
    }

    /// <summary>
    /// 
    /// </summary>
    /// <param name="_DelayTime"> 总执行多少时间（_isLoop true 时无效）</param>
    /// <param name="repeat"> 间隔</param>
    /// <param name="_action"> 事件</param>
    /// <param name="_isLoop"> 循环执行</param>
    /// <param name="isDelay"> 循环第一次是否延迟</param>
    public TimeItem(float _DelayTime, float repeat, Action _action, bool _isLoop = false, bool isDelay = false) : this()
    {
        DelayTime = _DelayTime;
        action = _action;
        repeatTime = repeat;
        _repeatTime = isDelay ? repeat : 0;
        timeType = TimeType.Repeat;
        IsUsed = false;
        isLoop = _isLoop;
    }

    public enum TimeType
    {
        Repeat,
        Normal,
    }
    public float DelayTime = 0f;
    public bool isLoop = false;
    public TimeType timeType = TimeType.Normal;
    public Action action;
    public float repeatTime = 1f;
    public float _repeatTime = 0f;
    public bool IsUsed = false;
}

public class GameHelper : MonoBehaviour
{
    private static GameHelper _instance;
    public static GameHelper Instance
    {
        get
        {

            return _instance;
        }
    }
    public string ActiveSenceName = "";
    GameHelpMessage gameMessage = new GameHelpMessage();

    private TimeItem ZeroCheakTime;

    public void ClearZeroCheakTime()
    {
        if (ZeroCheakTime != null)
        {
            RemoveInvoke(ZeroCheakTime);
            ZeroCheakTime = null;
        }
    }
    public void SetZeroCheakTime()
    {
        ClearZeroCheakTime();
        int seconds = (int)GameUtil.Instance.GetNowToTimeSecoud(1, 0, 00, 1);

        ZeroCheakTime = AddInvoke(new TimeItem(seconds, () =>
        {
            //JEvent.defaultEvent.Post(new GameEvent.Refresh_ZeroTime());
            UniGameEvent.Send(new GameEvent.Refresh_ZeroTime());
            SetZeroCheakTime(); // 重新设置
        }));
    }
    private List<TimeItem> TimeList = new List<TimeItem>();

    public TimeItem AddInvoke(TimeItem _add)
    {
        TimeList.Add(_add);
        return _add;
    }

    public void RemoveInvoke(TimeItem _Item)
    {
        TimeList.Remove(_Item);
    }

    List<TimeItem> itemsToRemove = new List<TimeItem>();

    private void Update()
    {
        GameM.Instance.SeverTime += Time.deltaTime;
        for (int i = 0; i < TimeList.Count; i++)
        {
            var _timeItem = TimeList[i];
            if (!_timeItem.IsUsed)
            {
                _timeItem.DelayTime -= Time.deltaTime;

                if (_timeItem.timeType == TimeItem.TimeType.Repeat)
                {
                    if (_timeItem._repeatTime <= 0)
                        _timeItem.action?.Invoke();

                    if (_timeItem._repeatTime <= 0)
                    {
                        _timeItem._repeatTime = _timeItem.repeatTime;
                    }
                    _timeItem._repeatTime -= Time.deltaTime;
                }
                else
                {
                    if (_timeItem.DelayTime <= 0)
                    {
                        _timeItem.action?.Invoke();
                    }
                }
                if (_timeItem.DelayTime < 0 && !_timeItem.isLoop)
                {
                    _timeItem.IsUsed = true;
                    itemsToRemove.Add(_timeItem);
                }
            }
        }
        if (itemsToRemove.Count > 0)
        {
            for (int i = 0; i < itemsToRemove.Count; i++)
            {
                RemoveInvoke(itemsToRemove[i]);
            }
            itemsToRemove.Clear();
        }

    }

    public void Awake()
    {
        JEvent.defaultEvent.Register(gameMessage);
        TimeList.Clear();
        if (_instance == null)
        {
            _instance = this;
            DontDestroyOnLoad(_instance);
        }
    }

    private void OnDestroy()
    {
        JEvent.defaultEvent.Unregister(gameMessage);
        _instance = null;
    }

    public void BackToLogin(bool isChange = false)
    {
        GameM.Instance.SeverTimeOff = 0;
        TCPClientMgr.Instance.ClearIpPool();

        if (isChange)
        {
            PlayerPrefsExtend.DeleteKey(SDKCtr.SaveKey.WX_REFRESH_TOKEN_TIME.ToString());
            PlayerPrefsExtend.DeleteKey(SDKCtr.SaveKey.WX_OPEN_ID.ToString());
            PlayerPrefsExtend.DeleteKey(SDKCtr.SaveKey.WX_ACCESS_TOKEN.ToString());
            PlayerPrefsExtend.DeleteKey(SDKCtr.SaveKey.WX_REFRESH_TOKEN.ToString());
        }

        JEvent.defaultEvent.Post(new SenceChange()
        {
            scencePath = GameConst.LoginPath,
        });
    }

    public async UniTask SetRoleAsync(HotUpdateScripts.ProtoM.JobType jobType, SexType sexType, Action<GameObject> action)
    {
        string poolKey = "MainRole" + jobType + sexType;
        if (!PoolHelper.Instance.CheckPool(poolKey))
        {
            ResouseHelper.Instance.LoadGameObj(GameConst.GetModelByJobAndSex(jobType, sexType), go =>
           {
               PoolHelper.Instance.AddPool(poolKey, go);
           }).Forget();
        }
        await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
        GameObject _mainRole = PoolHelper.Instance.GetPoolObj(poolKey);
        _mainRole.gameObject.name = poolKey;
        action?.Invoke(_mainRole);
    }

    public async UniTask SetAsyncObj(string poolKey, string path, Action<GameObject> action, bool isUsedPath = false)
    {
        if (!PoolHelper.Instance.CheckPool(poolKey))
        {
            await ResouseHelper.Instance.LoadGameObj((isUsedPath ? path : GameConst.GetModelPath(path)), go =>
            {
                PoolHelper.Instance.AddPool(poolKey, go);
            });
            //var baseObj = new JPrefab((isUsedPath ? path: GameConst.GetModelPath(path)), true);
            //await baseObj.WaitForAsyncLoading();
            //PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
        }
        await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
        GameObject _mainRole = PoolHelper.Instance.GetPoolObj(poolKey);
        _mainRole.gameObject.name = poolKey;
        action?.Invoke(_mainRole);
    }

    public void ClearObjToPool(Transform root)
    {
        if (root.childCount > 0)
        {
            for (int i = 0; i < root.childCount; i++)
            {
                PoolHelper.Instance.RecycObj(root.GetChild(i).gameObject);
            }
        }
    }

    //延迟下一帧
    public void SetDeayFrame(Action acation)
    {
        CoroutineMgr.Instance.StartCoroutine(SetCoroutine(acation));
    }

    private System.Collections.IEnumerator SetCoroutine(Action action)
    {
        yield return new WaitForEndOfFrame();
        action?.Invoke();
    }

    //延迟几秒
    public Coroutine SetDeayTimeDo(Action acation, float time)
    {
        return CoroutineMgr.Instance.StartCoroutine(SetTimeCoroutine(acation, time));
    }

    private System.Collections.IEnumerator SetTimeCoroutine(Action action, float time)
    {
        yield return new WaitForSeconds(time);
        action?.Invoke();
    }

    public void StarGame(int buattleId)
    {
        CFGItem_Battle cfg = CFG_Battle.GetById(buattleId);
        if (cfg.Powerlimit > RoleM.Instance.TotalCombat)
        {
            UITool.ShowTips("当前战力不足");
            return;
        }

        if (cfg.Levellimit > RoleM.Instance.Level)
        {
            UITool.ShowTips("当前等级不足");
            return;
        }

        if (cfg.Copytype == (int)CFG_Room.RoomType.Boss && FightPlayM.Instance.GetInfo().boss_map.ChallengeTime <= 0)
        {
            UITool.ShowTips("挑战次数不足");
            return;
        }
        TcpC.GetGameJoin(CFG_Room.GetUIDByRoomType((CFG_Room.RoomType)cfg.Copytype), buattleId);

        // 数数上报
        ShuShuSend((CFG_Room.RoomType)cfg.Copytype);
    }

    public void ShuShuSend(CFG_Room.RoomType roomType)
    {
        if (roomType == CFG_Room.RoomType.Fight)
        {
            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "玩法",
                component = "挑战",
                position = "战役",
                event_click = $"开始挑战",
            });
        }
        else if (roomType == CFG_Room.RoomType.SLTower)
        {
            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "玩法",
                component = "挑战",
                position = "试炼塔",
                event_click = $"开始挑战",
            });
        }
        else if (roomType == CFG_Room.RoomType.Relic)
        {
            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "玩法",
                component = "挑战",
                position = "遗迹",
                event_click = $"开始挑战",
            });
        }
        else if (roomType == CFG_Room.RoomType.Boss)
        {
            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "玩法",
                component = "挑战",
                position = "首领",
                event_click = $"挑战首领",
            });
        }
        else if (roomType == CFG_Room.RoomType.YZDD)
        {
            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "玩法",
                component = "竞技",
                position = "一站到底",
                event_click = $"开始挑战",
            });
        }
        else if (roomType == CFG_Room.RoomType.WorldBoss)
        {
            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "玩法",
                component = "竞技",
                position = "世界Boss",
                event_click = $"开始挑战",
            });
        }
    }

    void OnApplicationFocus(bool focus)
    {
        UniGameEvent.Send(new GameEvent.OnApplicationFocus { focus = focus });
    }

    void OnApplicationPause(bool pause)
    {
        UniGameEvent.Send(new GameEvent.OnApplicationPause { pause = pause });
    }

    #region 登录 TCP
    private static int ConnectCount = 0;
    /// <summary>
    /// 开启TCP链接
    /// </summary>
    public void OnConnecdTCP()
    {
        ConnectCount++;
        TCPClientMgr.Instance.AddIpToPool(GameM.Instance.SeverInfoDic[GameM.Instance.SeverDicIdx][GameM.Instance.SeverIdx].websocket, "");
        TCPClientMgr.Instance.OnSocketConnectResult = null;
        TCPClientMgr.Instance.OnSocketConnectResult += (_t) =>
        {
            Debug.Log("OnSocketConnectResult-----");
            OnLoginToSever(_t);
        };
        TCPClientMgr.Instance.OnSocketHasError = null;
        TCPClientMgr.Instance.OnSocketHasError += (_t) =>
        {
            Debug.Log("OnSocketHasError-----");
            UITool.Close<UILoading>();
            TCPClientMgr.Instance.ClearIpPool();
            bool isLogin = SceneManager.GetActiveScene().name == "Login";
            Action yesFun = delegate { OnConnecdTCP(); };

            Action NoFun = delegate
            {
                ConnectCount = 0;
                TCPClientMgr.Instance.ClearIpPool();
                if (isLogin)
                {
                    UITool.Close<UISelectRole>();
                    UITool.Close<UICreatRole>();
                    UITool.Open<UILogin>();
                }
                else
                {
                    JEvent.defaultEvent.Post(new SenceChange()
                    {
                        scencePath = GameConst.LoginPath,
                    });
                }
            };

            if (GameM.Instance.isSeverClose)
            {
                TCPClientMgr.Instance.OnSocketHasError = null;
                return;
            }
            if (ConnectCount > 3)
            {
                var info = new UIMessageBoxOne.CallBack
                {
                    showStr = "网络异常，请重新登录",
                    cType = 0,
                    endFun = NoFun,
                };
                UITool.Open<UIMessageBoxOne>(null, info);
            }
            else
            {
                UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack { ShowStr = "已断开连接，请重新连接", NoFun = NoFun, YesFun = yesFun });
            }

        };
        //UITool.Open<UILoading>(null, new UILoading.Loading { _time = 0, _str = "同步中..." });
        TCPClientMgr.Instance.ConnectSocket();
    }

    private bool _IsConnect = false;
    /// <summary>
    /// 检测链接信息
    /// </summary>
    /// <param name="group"></param>
    private void OnLoginToSever(WebClientGroup group)
    {
        if (group.isConnect)
        {
            _IsConnect = false;
            var LoginInfo = new LoginToGameServerRequest
            {
                Name = RoleM.Instance.Nickname,
                Uid = RoleM.Instance.Uid,
                //HeadUrl = RoleM.Instance.HeadUrl,
                Token = RoleM.Instance.Token,
            };
            TCPClientMgr.Instance.SendMessage((short)CmdId.C2S_Login, MType.Request, TCPClientMgr.ProtoSerialize(LoginInfo));
        }
        else
        {
            UITool.Close<UILoading>();
            TCPClientMgr.Instance.ClearIpPool();
            if (!_IsConnect)
            {
                _IsConnect = true;
                var callBack = new UIMessageBox.CallBack { ShowStr = "服务器连接失败，请重试", YesFun = OnConnecdTCP };
                if (SceneManager.GetActiveScene().name != "Login")
                {
                    callBack.NoFun = delegate { BackToLogin(); };
                }
                UITool.Open<UIMessageBox>(null, callBack);
            }
            else
            {
                _IsConnect = false;
                UITool.Open<UIMessageBoxOne>(null, new UIMessageBoxOne.CallBack
                {
                    showStr = "服务器连接失败，请稍候再试",
                    endFun = delegate { BackToLogin(); },
                });
            }
        }
    }
    #endregion

}

[Subscriber(ThreadMode.Main)]
public class GameHelpMessage
{
    private async UniTask OnSenceChange(SenceChange data)
    {
        GameHelper.Instance.ActiveSenceName = data.scencePath;
        string sencePath = $"Assets/HotUpdateResources/{data.scencePath}.unity";
        if (data._senceType == 1)
        {
            sencePath = $"Assets/{data.scencePath}.unity";
        }
        AudioManager.Instance.StopBgMusic();
        //await UniTask.DelayFrame(8);
        await AssetMgr.LoadSceneAsync(sencePath);
        //AudioManager.Instance.gameObject.SetActive(true);
        data.CallBack?.Invoke();
        //if (data.scencePath == GameConst.NorFightPath || data.scencePath == GameConst.TurnBasedFightPath)
        //AudioManager.Instance.PlayBgMusic(AudioKey.Sound.BGM_Fight);
        // else 

        if (data.scencePath == GameConst.HallPath)
            AudioManager.Instance.PlayBgMusic(AudioKey.Sound.BG_HallRest);
    }
}

/// <summary>
/// 游戏场景跳转
/// </summary>
public struct SenceChange
{
    public string scencePath;
    public Action CallBack;
    public int _senceType;
}
