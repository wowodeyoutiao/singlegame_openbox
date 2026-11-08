using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.UI;
using JEngine.Helper;
using HotUpdateScripts.ProtoM;
using SuperScrollView;
using System;
using pb;
using System.Linq;
using JEngine.Event;

public class UIVIP : UIBase<UIModelVIP, UIViewVIP>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    private List<CFGItem_Store> StoreList;
    private List<TaskInfo> taskInfos;
    private List<TaskInfo> curLevelInfos;

    public override void OnInit()
    {
        StoreList = CFG_Store.GetListByStoreType(CFG_Store.StoreType.VIPGift);

        GetView().btn_VipAdd.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            selectedVip++;
            if (selectedVip > CFG_VIP.GetMaxVipLevel())
                selectedVip = CFG_VIP.GetMaxVipLevel();

            SetVipShow();
        });

        GetView().btn_VipJian.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            selectedVip--;
            if (selectedVip < 1)
                selectedVip = 1;

            SetVipShow();
        });

        GetView().InfiniteScroll_Mission.InitListView(0, OnFillItemMission);

        GetView().btn_VipGiftBuy.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (IsVIPUnlock(selectedVip))
            {
                UITool.ShowTips("您当前的VIP等级不足，请提升VIP等级后再来购买！");
                return;
            }

            if (!GameUtil.Instance.ComCoinToCheck(StoreList[selectedVip - 1].ID, () => UITool.Close<UIVIP>()))
            {
                return;
            }

            GameUtil.Instance.GameToBuy(StoreList[selectedVip - 1].ID, () =>
            {
                //RedPointManager.Instance.RPS.Dispatch(RedPointConst.VIP_Buy.ToString(), TaskM.Instance.VIPCanBuyGiftCount());
                SetVipShow();
            });

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "首页",
                component = "VIP按钮",
                position = "VIP界面",
                event_click = "VIP礼包购买",
            });

        });

        GetView().InfiniteScroll_Des.InitListView(0, OnFillItemDesc);

        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIVIP>();
        });
    }


    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Sync_UserBase>(Event_UserBase);

        RefreshUI();
    }

    //public void OnBindRed()
    //{
    //    RedPointManager.Instance.RPS.Register(RedPointConst.VIP_Buy.ToString(), (node) =>
    //    {
    //        GetView().btn_VipGiftBuy.transform.Find("dot").gameObject.SetActive(node.pointCount > 0);
    //    });
    //}

    private void RefreshUI()
    {
        GetVIPMissionList();
        ShowCurrentVIPLevel();

        selectedVip = RoleM.Instance.Vip < CFG_VIP.GetMaxVipLevel() ? RoleM.Instance.Vip + 1 : CFG_VIP.GetMaxVipLevel();
        SetVipShow();
    }

    [Subscriber(ThreadMode.Main)]
    private void Task_Refresh(GameEvent.Refresh_Task _e)
    {
        if (isNullGo)
            return;

        GetVIPMissionList();
        ShowCurrentVIPLevel();
        SetVipShow();
    }

    [Subscriber(ThreadMode.Main)]
    private void OnSync_HomeOpenBox(GameEvent.Sync_HomeOpenBox _e)
    {
        if (isNullGo)
            return;

        if (_e.data.Code == 0 && ExistOpenBoxMission())
            RefreshMissionList();
    }

    //[Subscriber(ThreadMode.Main)]
    private void Event_UserBase(GameEvent.Sync_UserBase _e)
    {
        if (isNullGo)
            return;

        ShowCurrentVIPLevel();
    }

    private void GetVIPMissionList()
    {
        taskInfos = TaskM.Instance.GetTasksByClass(TaskClass.VIP);
    }

    // vip还未解锁
    private bool IsVIPUnlock(int vipLv)
    {
        return vipLv > RoleM.Instance.Vip;
    }

    private LoopListViewItem2 OnFillItemMission(LoopListView2 loopListView, int idx)
    {
        LoopListViewItem2 item = loopListView.NewListViewItem("Item");

        TaskInfo curInfo = curLevelInfos[idx];
        int taskId = curInfo.CfgId;

        GameObject go = item.gameObject;
        go.name = "Item" + idx;
        JUIExtensions.txt(go.transform.Find("txt_condition")).text = CFG_Task.StrTaskDes(taskId);
        (int amout, int targetnum) = CFG_Task.GetTaskPro(taskId, curInfo.Target.Amount, curInfo.Target.TargetNum);
        string colorHex = amout < targetnum ? "C7351D" : "43A019";
        JUIExtensions.txt(go.transform.Find("txt_progress")).text = $"<size=14>进度 </size><color=#{colorHex}>{GameUtil.Instance.FormatNumWithUnit(amout)}</color>/{GameUtil.Instance.FormatNumWithUnit(targetnum)}";
        Transform awardRoot = go.transform.Find("Awards");

        CFGItem_Task cfg_task = CFG_Task.GetById(taskId);

        go.transform.Find("txt_geted").gameObject.SetActive(curInfo.Status == TaskStatus.REWARDED);
        Transform trans_btn_go = go.transform.Find("btn_go");
        trans_btn_go.gameObject.SetActive(curInfo.Status == TaskStatus.IN_PROGRESS);
        JUIExtensions.btn(trans_btn_go).onClick.RemoveAllListeners();
        JUIExtensions.btn(trans_btn_go).onClick.AddListener(() =>
        {
            // 跳转到对应的任务界面
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            CFG_Task.JumpUI(taskId, () =>
            {
                UITool.Close<UIVIP>();
            });
        });

        Transform trans_btn_get = go.transform.Find("btn_get");
        trans_btn_get.gameObject.SetActive(curInfo.Status == TaskStatus.COMPLETED);
        JUIExtensions.btn(trans_btn_get).onClick.RemoveAllListeners();
        JUIExtensions.btn(trans_btn_get).onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            if (amout >= targetnum)
            {
                OnTaskGet(taskId);

                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "首页",
                    component = "VIP按钮",
                    position = "VIP界面",
                    event_click = $"任务完成领取奖励",
                });
            }
            else
            {
                UITool.ShowTips("任务未完成，无法领取奖励！");
            }
        });


        for (int i = 0; i < awardRoot.childCount; i++)
        {
            Transform itemAward = awardRoot.GetChild(i);
            bool isShow = i < cfg_task.Reward.Length;
            itemAward.gameObject.SetActive(isShow);
            if (isShow)
            {
                int cfgId = int.Parse(cfg_task.Reward[i].ElementAt(0).Key);
                int num = cfg_task.Reward[i].ElementAt(0).Value;
                GameUtil.Instance.SetItemIcon(cfgId, JUIExtensions.img(itemAward.Find("icon")), JUIExtensions.img(itemAward));
                JUIExtensions.txt(itemAward.Find("txt_num")).text = GameUtil.Instance.FormatNumWithUnit(num);
                JUIExtensions.btn(itemAward).onClick.RemoveAllListeners();
                JUIExtensions.btn(itemAward).onClick.AddListener(() =>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfgId });
                });
            }
        }

        return item;
    }

    private LoopListViewItem2 OnFillItemDesc(LoopListView2 loopListView, int idx)
    {
        LoopListViewItem2 item = loopListView.NewListViewItem("Item");
        GameObject go = item.gameObject;
        JUIExtensions.txt(go.transform.Find("txt")).text = privilegedDescription[idx];
        return item;
    }

    private int selectedVip;
    private string[] privilegedDescription;
    private CFGItem_VIP curVIP_CFG;

    private void ShowCurrentVIPLevel()
    {
        GetView().txt_CurrentVipLevel.text = $"VIP{RoleM.Instance.Vip}";
        int nextLv = RoleM.Instance.Vip < CFG_VIP.GetMaxVipLevel() ? RoleM.Instance.Vip + 1 : CFG_VIP.GetMaxVipLevel();
        GetView().txt_NextVipLevel.text = $"VIP{nextLv}";
        (long currentExp, int currentLevelEXPNeed, float progress) = GetVIPProgress(nextLv);
        GetView().txt_Progress.text = currentExp + "/" + currentLevelEXPNeed;
        GetView().slider_VIP.value = progress;
    }

    private void SetVipShow()
    {
        curVIP_CFG = CFG_VIP.GetById(selectedVip);
        GetView().txt_Title.text = $"VIP{selectedVip}";
        // 任务详情
        GetView().txt_VipMission.text = $"VIP {selectedVip} 任务";
        GetCurrentLevelVIPTasks();
        RefreshMissionList();

        // 任务未解锁
        GetView().vip_MissionUnlockBG.SetActive(IsVIPUnlock(selectedVip - 1));
        if (IsVIPUnlock(selectedVip - 1))
            GetView().txt_VipMissionUnlockTip.text = $"VIP {selectedVip - 1} 解锁后开启";

        // 切换等级希显示按钮
        GetView().btn_VipJian.gameObject.SetActive(selectedVip > 1);
        GetView().btn_VipAdd.gameObject.SetActive(selectedVip < CFG_VIP.GetMaxVipLevel());
        // 特权说明
        GetView().txt_vipPrivilege.text = $"VIP {selectedVip} 特权";
        privilegedDescription = curVIP_CFG.Pdisplay.Split('|');
        GetView().InfiniteScroll_Des.SetListItemCount(privilegedDescription.Length);
        GetView().InfiniteScroll_Des.RefreshAllShownItem();
        // 礼包信息
        GetView().txt_vipGift.text = $"VIP {selectedVip} 礼包";
        GetView().txt_Discount.text = curVIP_CFG.Discount.ToString() + "折！";

        var storeInfo = StoreList[selectedVip - 1];
        for (int i = 0; i < GetView().root_Awards.childCount; i++)
        {
            GetView().root_Awards.GetChild(i).gameObject.SetActive(i < storeInfo.Solditems.Length);

            if (i < storeInfo.Solditems.Length)
            {
                Transform item = GetView().root_Awards.GetChild(i);
                GameUtil.Instance.SetItemIcon(storeInfo.Solditems[i].itemID, JUIExtensions.img(item.Find("icon")), JUIExtensions.img(item));
                JUIExtensions.txt(GetView().root_Awards.GetChild(i).Find("txt_num")).text = GameUtil.Instance.FormatNumWithUnit(storeInfo.Solditems[i].Num);

                int cfgId = storeInfo.Solditems[i].itemID;
                JUIExtensions.btn(GetView().root_Awards.GetChild(i)).onClick.RemoveAllListeners();
                JUIExtensions.btn(GetView().root_Awards.GetChild(i)).onClick.AddListener(() =>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = cfgId });
                });
            }
        }

        // 礼包购买状态
        bool showBuyButton = !IsVIPUnlock(selectedVip) && StoreM.Instance.GetBuyNumByID(storeInfo.ID) < storeInfo.Limite.limitnum;
        if (showBuyButton)
        {
            GameUtil.Instance.SetItemIcon(storeInfo.Limite.needmoney, GetView().img_GiftBuy);
            string colorHex = BagM.Instance.GetNumByCfgID(storeInfo.Limite.needmoney) < storeInfo.Limite.Num ? "C7351D" : "36901E";
            GetView().txt_GiftBuy.text = $"<color=#{colorHex}>{GameUtil.Instance.FormatNumWithUnit(storeInfo.Limite.Num)}</color>";
        }
        GetView().btn_VipGiftBuy.gameObject.SetActive(showBuyButton);
        GetView().buy_Consume.SetActive(showBuyButton);
        GetView().txt_HaveBought.gameObject.SetActive(StoreM.Instance.GetBuyNumByID(storeInfo.ID) >= storeInfo.Limite.limitnum);
        GetView().txt_VIPUnlockTip.gameObject.SetActive(IsVIPUnlock(selectedVip));
        if (IsVIPUnlock(selectedVip))
            GetView().txt_VIPUnlockTip.text = $"VIP {selectedVip} 未解锁";
    }

    private void RefreshMissionList()
    {
        int curTaskCount = IsVIPUnlock(selectedVip - 1) ? 0 : curLevelInfos.Count;
        GetView().InfiniteScroll_Mission.SetListItemCount(curTaskCount);
        GetView().InfiniteScroll_Mission.RefreshAllShownItem();
    }

    private void GetCurrentLevelVIPTasks()
    {
        curLevelInfos = taskInfos.FindAll(v =>
        {
            return CFG_Task.GetById(v.CfgId).VIPLevel == selectedVip;
        });
    }

    /// <summary>
    /// 任务列表是否存在开宝箱任务
    /// </summary>
    /// <returns>bool</returns>
    private bool ExistOpenBoxMission()
    {
        return curLevelInfos != null && curLevelInfos.Count > 0 && curLevelInfos.Exists(v => v.Target.Type == TargetType.TOTAL_BOX_OPENED);
    }

    private (long, int, float) GetVIPProgress(int nextLevel)
    {
        int currentLevelEXPNeed = CFG_VIP.GetById(nextLevel).NeedEXP;
        float progress = (float)RoleM.Instance.VipExp / currentLevelEXPNeed;
        return (RoleM.Instance.VipExp, currentLevelEXPNeed, progress);
    }

    private void OnTaskGet(int taskId)
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        TcpC.C2S_TaskComplate(taskId);
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Sync_UserBase>(Event_UserBase);
        //RedPointManager.Instance.RPS.Register(RedPointConst.VIP_Buy.ToString(), null);
    }

    public override void OnDestroy()
    {

    }
}