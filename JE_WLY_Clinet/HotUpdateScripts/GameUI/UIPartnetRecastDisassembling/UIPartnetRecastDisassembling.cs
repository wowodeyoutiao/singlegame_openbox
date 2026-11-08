using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts;
using JEngine.Helper;
using System;
using HotUpdateScripts.ProtoM;
using SuperScrollView;
using JEngine.Event;
using pb;
using BestHTTP.SecureProtocol.Org.BouncyCastle.Asn1.Ocsp;
using JEngine.UI;
using ILRuntime.Runtime.Debugger.Expressions;
using System.Linq;
using Spine;


public class UIPartnetRecastDisassembling : UIBase<UIModelPartnetRecastDisassembling, UIViewPartnetRecastDisassembling>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(()=> {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIPartnetRecastDisassembling>();
        });

        GetView().btn_Help.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            int id = partnerOptionType == PartnerOptionType.Recast ? 38 : 39;
            UITool.Open<UIComRule>(null, id);
        });

        GetView().loopList_Item.InitListView(0, OnItemFill);

        GetView().btn_Recast.onClick.AddListener(OnRecastPartner);
        GetView().btn_Disassembling.onClick.AddListener(OnDisassemblingPartner);
    }

    public enum PartnerOptionType
    {
        /// <summary>
        /// 伙伴重铸
        /// </summary>
        Recast,

        /// <summary>
        /// 伙伴分解
        /// </summary>
        Disassembling
    }

    PartnerOptionType partnerOptionType;

    public long partnerUid;
    public int slotIndex;
    public PartnerInfo partnerInfo;
    CFGItem_Partner cfgItem_Partner;
    CFGItem_Item cfgItem_Item;

    public override void OnEnter(params object[] args)
    {
        //JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Refrsh_Common>(OnRefresh_CommonResponse);

        partnerUid = (long)args[0];
        partnerInfo = PartnerM.Instance.GetOnePartnerInfoByUid(partnerUid);
        slotIndex = (int)args[1];

        if (args.Length > 2)
        {
            partnerOptionType = (PartnerOptionType)args[2];
        }
        RefreshUI();
    }

    private LoopListViewItem2 OnItemFill(LoopListView2 list, int index)
    {
        var item = list.NewListViewItem("Item");
        Transform t = item.transform;
        GameUtil.Instance.SetItemIcon(items[index].id, JUIExtensions.img(t.Find("icon")), JUIExtensions.img(t));
        JUIExtensions.txt(t.Find("txt_num")).text = items[index].num.ToString();

        int idx = index;
        JUIExtensions.btn(t).onClick.RemoveAllListeners();
        JUIExtensions.btn(t).onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = items[idx].id });
        });

        return item;
    }

    private void RefreshUI()
    {
        cfgItem_Partner = CFG_Partner.GetById(partnerInfo.CfgId);
        cfgItem_Item = CFG_Item.GetById(cfgItem_Partner.ItemID);

        if (partnerOptionType == PartnerOptionType.Recast)
        {
            GetView().txt_Title.text = "伙伴重铸";
            GetView().txt_Desc.text = "重铸后可获得";
            GetView().toggle_RecastFree.isOn = true;
            GetView().txt_FreeRecastDesc.text = $"免费重铸  (材料折损{CFG_XMLConfig.xmlConfig.Partener.forgeRate}%)";
            GetView().txt_PayRecastDesc.text = $"付费重铸  (材料无折损)";
            int consumeId = CFG_XMLConfig.xmlConfig.PartenerPayRecastingDic[(int)cfgItem_Item.Quality].id;
            int consumeNum = CFG_XMLConfig.xmlConfig.PartenerPayRecastingDic[(int)cfgItem_Item.Quality].num;
            GameUtil.Instance.SetItemIcon(consumeId, GetView().img_PayRecastConsume);
            bool isEnough = GameUtil.Instance.ComCoinToCheck(consumeId, consumeNum, null, false);
            string consumeStr = isEnough ? $"<color=#6BCB47>{GameUtil.Instance.FormatNumWithUnit(consumeNum)}</color>" : $"<color=#DA4E38>{GameUtil.Instance.FormatNumWithUnit(consumeNum)}</color>";
            GetView().txt_PayRecastConsume.text = consumeStr;
        }
        else
        {
            GetView().txt_Title.text = "伙伴分解";
            GetView().txt_Desc.text = "分解后可获得";
        }


        GetView().txt_Name.text = GameUtil.Instance.GetItemName(cfgItem_Partner.ItemID, -1);
        JUIExtensions.outline(GetView().txt_Name.transform).effectColor = GameUtil.Instance.GetQualityOutlineColor((int)cfgItem_Item.Quality);
        GameUtil.Instance.SetItemIcon(cfgItem_Partner.ItemID, GetView().img_Icon, GetView().img_IconBG);
        GameUtil.Instance.SetJobIcon(CFG_Partner.GetById(partnerInfo.CfgId).job % 10, GetView().img_Job);

        GetView().root_Toggle.gameObject.SetActive(partnerOptionType == PartnerOptionType.Recast);
        //bool showRecastPay = partnerOptionType == PartnerOptionType.Recast && partnerInfo.Level >= PartnerM.payRecastLevel;
        //GetView().toggle_RecastPay.gameObject.SetActive(showRecastPay);
        //GetView().txt_PayRecastConsume.gameObject.SetActive(showRecastPay);
        GetView().toggle_RecastFree.onValueChanged.AddListener(isOn =>
        {
            if (isOn)
                ShowAwards();
        });
        GetView().toggle_RecastPay.onValueChanged.AddListener(isOn =>
        {
            if (isOn)
                ShowAwards();
        });

        GetView().btn_Recast.gameObject.SetActive(partnerOptionType == PartnerOptionType.Recast);
        GetView().btn_Disassembling.gameObject.SetActive(partnerOptionType == PartnerOptionType.Disassembling);

        ShowAwards();
    }

    public class Item
    {
        public Item() { }

        public Item(int id, int num) 
        {
            this.id = id;
            this.num = num;
        }

        public int id;
        public int num;
    }

    public List<Item> items = new List<Item>();

    private void ShowAwards()
    {
        items.Clear();
        var debrisItems = GetDebrisItems();

        if (debrisItems.Count > 0)
            items.AddRange(debrisItems);

        var expItems = GetEXPItems();
        if (expItems.Count > 0)
            items.AddRange(expItems);

        GetView().loopList_Item.SetListItemCount(items.Count);
        GetView().loopList_Item.RefreshAllShownItem();
    }

    private List<Item> GetDebrisItems()
    {
        List<Item> debrisItems = new List<Item>();
        var debris = cfgItem_Partner.GetStarAllConsumeFromLv0(partnerInfo.Star);

        if (debris.Length > 0)
        {
            // 分解
            if (partnerOptionType == PartnerOptionType.Disassembling)
            {
                int id1 = CFG_XMLConfig.xmlConfig.Partener.decompose.id;
                int num1 = Mathf.FloorToInt(debris[0].value * (CFG_XMLConfig.xmlConfig.Partener.decompose.rate / 100f));
                debrisItems.Add(new Item(id1, num1));

                int id2 = debris[1].key;
                int num2 = Mathf.FloorToInt(debris[1].value * (CFG_XMLConfig.xmlConfig.Partener.decompose.rate / 100f));
                debrisItems.Add(new Item(id2, num2));
            }
            // 免费重铸
            else if (GetView().toggle_RecastFree.isOn)
            {
                int id1 = debris[0].key;
                int num1 = Mathf.FloorToInt(debris[0].value * (CFG_XMLConfig.xmlConfig.Partener.forgeRate / 100f));
                debrisItems.Add(new Item(id1, num1));

                int id2 = debris[1].key;
                int num2 = Mathf.FloorToInt(debris[1].value * (CFG_XMLConfig.xmlConfig.Partener.forgeRate / 100f));
                debrisItems.Add(new Item(id2, num2));
            }
            // 付费重铸
            else
            {
                int id1 = debris[0].key;
                int num1 = Mathf.FloorToInt(debris[0].value);
                debrisItems.Add(new Item(id1, num1));

                int id2 = debris[1].key;
                int num2 = Mathf.FloorToInt(debris[1].value);
                debrisItems.Add(new Item(id2, num2));
            }
        }

        return debrisItems;
    }

    private List<Item> GetEXPItems()
    {
        List<Item> expItems = new List<Item>();
        long allExpNeed = CFG_Role.GetAllNeedExpToUpLevel(cfgItem_Partner.job, partnerInfo.Level);

        // 免费重铸（经验折算）,分解或付费重铸经验不折算
        if (partnerOptionType == PartnerOptionType.Recast && GetView().toggle_RecastFree.isOn)
        {
            allExpNeed = Mathf.FloorToInt(allExpNeed * (CFG_XMLConfig.xmlConfig.Partener.forgeRate / 100f));
        }

       var expBeads = CFG_XMLConfig.xmlConfig.PartenerEXPBeadDic;

        for (int i = expBeads.Count - 1; i >= 0; i--)
        {
            var item = expBeads.ElementAt(i);
            if (allExpNeed >= item.Value)
            {
                int count = Mathf.FloorToInt((float)allExpNeed / item.Value);
                expItems.Add(new Item(item.Key, count));
                allExpNeed = allExpNeed % item.Value;
            }
        }

        return expItems;
    }

    public int EquipCountInPartner()
    {
        int count = 0;
        for (int i = 0; i < 12; i++)
        {
            var itemInfo = PartnerM.Instance.GetOnePartnerOneEquip(slotIndex, i);
            if (itemInfo != null)
                count++;
        }
        return count;
    }

    private void OnRecastPartner()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        if (!BagM.Instance.CheckEquipVacancy(EquipCountInPartner()))
        {
            UITool.ShowTips("包裹已满，无法重铸。");
            return;
        }

        if (GetView().toggle_RecastPay.isOn)
        {
            int consumeId = CFG_XMLConfig.xmlConfig.PartenerPayRecastingDic[(int)cfgItem_Item.Quality].id;
            int consumeNum = CFG_XMLConfig.xmlConfig.PartenerPayRecastingDic[(int)cfgItem_Item.Quality].num;

            bool isEnough = GameUtil.Instance.ComCoinToCheck(consumeId, consumeNum, () => { UITool.Close<UIPartnetRecastDisassembling>(); }, true);
            if (!isEnough)
                return;
        }

        TcpC.C2S_PartnerReforge(partnerUid, GetView().toggle_RecastFree.isOn);

        string pay = GetView().toggle_RecastFree.isOn ? "免费重铸" : "付费重铸";
        HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
        {
            page = "角色",
            component = "伙伴重铸",
            position = "",
            event_click = pay,
        });
    }

    private void OnDisassemblingPartner()
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        if (!BagM.Instance.CheckEquipVacancy(EquipCountInPartner()))
        {
            UITool.ShowTips("包裹已满，无法分解。");
            return;
        }

        TcpC.C2S_PartnerDisassemble(partnerUid);

        HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
        {
            page = "角色",
            component = "伙伴分解",
            position = "",
            event_click = "分解",
        });
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnRefresh_CommonResponse(GameEvent.Refrsh_Common _e)
    {
        if (isNullGo) return;

        // 伙伴重铸
        if (_e.code == 0 && _e.cid == CmdId.C2S_PartnerReforge)
        {
            UITool.Close<UIPartnetRecastDisassembling>();
        }

        // 伙伴分解
        if (_e.code == 0 && _e.cid == CmdId.C2S_PartnerDisassemble)
        {
            UITool.Close<UIPartnetRecastDisassembling>();
        }
    }

    public override void OnExit()
    {
        //JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refrsh_Common>(OnRefresh_CommonResponse);
        items.Clear();
    }

    public override void OnDestroy()
    {

    }
}