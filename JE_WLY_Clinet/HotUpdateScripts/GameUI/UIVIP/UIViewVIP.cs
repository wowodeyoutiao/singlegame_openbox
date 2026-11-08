using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewVIP : MonoBehaviour,IUIView
{
    // 标题文本
    public Text txt_Title;
    // vip等级和vip升级进度
    public Text txt_CurrentVipLevel;
	public Text txt_NextVipLevel;
	public Slider slider_VIP;
	public Text txt_Progress;
    // vip任务
    public Text txt_VipMission;
    public Button btn_VipJian;
	public Button btn_VipAdd;
    public LoopListView2 InfiniteScroll_Mission;
    public GameObject vip_MissionUnlockBG;
    public Text txt_VipMissionUnlockTip;
    // vip 说明
    public Text txt_vipPrivilege;
	public LoopListView2 InfiniteScroll_Des;
    // 礼包
    public Text txt_vipGift;
	public Transform root_Awards;
	public GameObject buy_Consume;
	public Image img_GiftBuy;
	public Text txt_GiftBuy;
	public Button btn_VipGiftBuy;
	public Text txt_HaveBought;
	public Text txt_VIPUnlockTip;
	public Text txt_PurchaseLimit;
	public Text txt_Discount;

    public Button btn_Close;

    public void Awake()
    {
        
    }

	public void Init(GameObject go)
	{

	}

	public void OnDestroy()
	{
		 UITool.RemoveUI<UIVIP>();
	}
}