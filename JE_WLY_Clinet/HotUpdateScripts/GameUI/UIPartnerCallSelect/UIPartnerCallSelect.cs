using Cysharp.Threading.Tasks;
using HotUpdateScripts;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using JEngine.UI;
using MetaFramework.UI;
using pb;
using Spine.Unity;
using System;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class UIPartnerCallSelect : UIBase<UIModelPartnerCallSelect, UIViewPartnerCallSelect>
{
    private UIViewPartnerCallSelect selectView;

    private List<Transform> itemTransList = new List<Transform>();
    private List<GameObject> selGoList = new List<GameObject>();
    private List<CFGItem_Partner> partnerCfgs = new List<CFGItem_Partner>();
    private int selItemId;
    private int oldSelItemId;
    private Action<int> callBack;

    private int[] priceInfo = new int[2];

    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    private void ButtonAddListener(Button btn, Action act)
    {
        btn.onClick.RemoveAllListeners();
        btn.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            act?.Invoke();
        });
    }

    public override void OnInit()
    {
        selectView = GetView();
        UniGameEvent.AddEventListener<LotterySelectResponse>(Sever_LotterySelectResponse);

        ButtonAddListener(selectView.closeBtn, () => { UITool.Close<UIPartnerCallSelect>(); });
        ButtonAddListener(selectView.okBtn, OnClickOkButton);
    }

    private void OnClickOkButton()
    {
        if (selItemId < 0)
        {
            UITool.ShowTips("请选择一个伙伴!");
            return;
        }

        if (priceInfo[0] > priceInfo[1])
        {
            UITool.ShowTips("道具不足，可通过运营活动获得!");
            return;
        }

        if (oldSelItemId == selItemId)
        {
            UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
            {
                ShowStr = $"您已经指定过该伙伴，是否继续消耗“心愿指定币”来重复指定？",
                YesFun = () =>
                {
                    TcpC.C2S_LotterySelect(selItemId);
                },
            });
            return;
        }

        TcpC.C2S_LotterySelect(selItemId);
    }

    private void Sever_LotterySelectResponse(LotterySelectResponse resq)
    {
        if (resq.Code != 0)
        {
            UITool.ShowTips(resq.Error);
            return;
        }

        callBack?.Invoke(selItemId);
        oldSelItemId = selItemId;
        UITool.ShowTips("选择成功！");
        RefreshPrice();
    }

    public override void OnEnter(params object[] args)
    {
        selItemId = (int)args[0];
        oldSelItemId = selItemId;
        callBack = (Action<int>)args[1];
        var num = CFG_XMLConfig.xmlConfig.GaCha.maxnum;
        selectView.titleText.text = $"<color=#845835>每当幸运值满{num}时刷新必出红色品质</color>\n<color=#ff5855>（指定刷新伙伴需要消耗道具，已获得的伙伴将会自动转换为对应数量碎片）</color>";

        LoadUI();

        RefreshPrice();
    }

    private void RefreshPrice()
    {
        var priceId = CFG_XMLConfig.xmlConfig.GaCha.id;
        var priceNum = CFG_XMLConfig.xmlConfig.GaCha.num;

        if (selItemId == 0)
            priceNum = 0;

        var bagNum = BagM.Instance.GetNumByCfgID(priceId);

        var colorStr = (bagNum >= priceNum ? "#209b23" : "#ff5855");
        selectView.priceText.text = $"<color={colorStr}>{bagNum}/{priceNum}</color>";

        priceInfo[0] = priceNum;
        priceInfo[1] = (int)bagNum;
    }

    private void LoadUI()
    {
        partnerCfgs = CFG_Partner.GetByQuality(7);
        if (partnerCfgs == null || partnerCfgs.Count == 0)
        {
            UITool.Close<UIPartnerCallSelect>();
            return;
        }
        selGoList.Clear();

        for (int i = 0; i < partnerCfgs.Count + 1; i++)
        {
            var index = i;
            Transform itemTrans;
            if (i < itemTransList.Count)
                itemTrans = itemTransList[i];
            else
            {
                itemTrans = selectView.GetItemTrans();
                itemTransList.Add(itemTrans);
            }
            itemTrans.localPosition = Vector3.zero;
            itemTrans.localScale = Vector3.one;
            itemTrans.gameObject.SetActive(true);
            
            // 背景
            var bgImg = JUIExtensions.img(itemTrans.Find("BgImg"));

            var ntureImg = JUIExtensions.img(itemTrans.Find("NatureImg"));
            ntureImg.enabled = false;

            var selImgGo = itemTrans.Find("SelButton/SelImg").gameObject;
            selGoList.Add(selImgGo);

            var selBtn = itemTrans.Find("SelButton").GetComponent<Button>();
            selBtn.onClick.RemoveAllListeners();
            selBtn.onClick.AddListener(() =>
            {
                selImgGo.SetActive(true);
                OnClickSelectAtIndex(index);
            });

            var freeGo = itemTrans.Find("FreeImg").gameObject;
            freeGo.SetActive(i == 0);
            if (i == 0)
            {
                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.PartnerCall, "bg_Summoning_dzk_7", bgImg);
                selImgGo.SetActive(selItemId == 0);
                JUIExtensions.txt(itemTrans.Find("NameText")).text = "随机";
                continue;
            }

            var partnerCfg = partnerCfgs[i - 1];
            var itemCfg = CFG_Item.GetById(partnerCfg.ItemID);
            selImgGo.SetActive(selItemId == partnerCfg.ItemID);

            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.PartnerCall, "bg_Summoning_dzk_" + itemCfg.Quality, bgImg);

            var bgBtn = itemTrans.Find("BgImg").GetComponent<Button>();
            bgBtn.onClick.RemoveAllListeners();
            bgBtn.onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = itemCfg.ID });
            });

            // 职业
            GameUtil.Instance.SetItemIcon(itemCfg.ID, null, null, ntureImg);

            // 角色
            if(itemTrans.Find("IconImg").childCount == 0)
            {
                string poolKey = "Partner_" + partnerCfg.ID;
                var path = GameConst.GetMonsterModelByID(poolKey);
                GameHelper.Instance.SetAsyncObj(poolKey, path, (role) =>
                {
                    if (isNullGo)
                        return;

                    role.transform.SetParent(itemTrans.Find("IconImg"));
                    role.name = "SpineItem";
                    role.transform.localScale = Vector3.one * 0.32f;
                    role.transform.localPosition = new Vector3(0, -26, 0);
                    role.SetActive(true);

                    var skeletonGraphic = role.GetComponentInChildren<SkeletonGraphic>();
                    SpineUtil.Instance.PlayLoop(skeletonGraphic, "stand");
                }).Forget();
            }

            // 名字
            JUIExtensions.txt(itemTrans.Find("NameText")).text = partnerCfg.Name;
        }
    }

    private void OnClickSelectAtIndex(int index)
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        for (int i = 0; i < selGoList.Count; i++)
            selGoList[i].SetActive(index == i);

        if (index != 0)
            selItemId = partnerCfgs[index - 1].ItemID;
        else
            selItemId = 0;

        RefreshPrice();
    }

    public override void OnExit()
    {

    }

    public override void OnDestroy()
    {
        OnExit();
        selectView = null;
        UniGameEvent.RemoveEventListener<LotterySelectResponse>(Sever_LotterySelectResponse);
    }
}