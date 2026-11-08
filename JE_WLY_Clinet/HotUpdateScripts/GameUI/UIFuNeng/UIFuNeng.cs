using DG.Tweening;

using HotUpdateScripts.Config;
using HotUpdateScripts.ProtoM;

using JEngine.Event;
using JEngine.Helper;

using MetaFramework.UI;

using Spine;

using System;
using System.Text;

using UnityEngine;

public class UIFuNeng : UIBase<UIModelFuNeng, UIViewFuNeng>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    
    public override void OnInit()
    {
        GetView().btn_Help.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIComRule>(null, 46);
        });

        GetView().btn_AddFuNeng.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIGiftBag>(null, ActivityType.Element);
        });

        GetView().btn_AddGold.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIShop>(null, CFG_Store.StoreType.Coupon);
        });

        GetView().btn_AddJinJie.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIGiftBag>(null, ActivityType.Element);
        });

        for (int i = 0; i < GetView().elements.Length; i++)
        {
            int pos = i + 1;
            GetView().elements[i].btn.onClick.AddListener(() =>
            {
                if (isAutoFuNeng)   // 自动赋能时不允许操作
                    return;

                SelectElement(pos);
            });
        }

        GetView().btn_FuNeng.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            if (!isCanOption) return;

            isAutoFuNeng = false;
            if (!IsEnoughConsumeFN(out int itemId))
            {
                isAutoFuNeng = false;
                UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
                {
                    ShowStr = $"当前持有的{GameUtil.Instance.GetItemName(itemId)}数量不足，是否前往获取？",
                    YesFun = () =>
                    {
                        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                        if((CFG_Item.ItemID)itemId == CFG_Item.ItemID.Gold)
                            UITool.Open<UIShop>(null, CFG_Store.StoreType.Coupon);
                        else
                            UITool.Open<UIGiftBag>(null, ActivityType.Element);
                    },
                });
                return;
            }

            isCanOption = false;
            OnFuNeng();

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "五行试炼界面",
                component = "赋能界面",
                position = "",
                event_click = $"单次赋能",
            });
        });

        GetView().btn_AutoFuNeng.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            if (isAutoFuNeng)
            {
                isAutoFuNeng = false;
                GetView().btn_FuNeng.gameObject.SetActive(true);
                GetView().txt_AutoFuNeng.text = "一键赋能";
                ClearAutoFuNengItem();

                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "五行试炼界面",
                    component = "赋能界面",
                    position = "",
                    event_click = $"取消自动赋能",
                });
            }
            else
            {
                if (!IsEnoughConsumeFN(out int itemId))
                {
                    isAutoFuNeng = false;
                    UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
                    {
                        ShowStr = $"当前持有的{GameUtil.Instance.GetItemName(itemId)}数量不足，是否前往获取？",
                        YesFun = () =>
                        {
                            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

                            if ((CFG_Item.ItemID)itemId == CFG_Item.ItemID.Gold)
                                UITool.Open<UIShop>(null, CFG_Store.StoreType.Coupon);
                            else
                                UITool.Open<UIGiftBag>(null, ActivityType.Element);
                        },
                    });
                    return;
                }

                isAutoFuNeng = true;
                GetView().btn_FuNeng.gameObject.SetActive(false);
                GetView().txt_AutoFuNeng.text = "停止赋能";
                OnFuNeng();

                HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
                {
                    page = "五行试炼界面",
                    component = "赋能界面",
                    position = "",
                    event_click = $"自动赋能",
                });
            }
        });

        GetView().btn_JinJie.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            if (!isCanOption) return;

            if (!IsEnoughConsumeJinJie(out int itemId))
            {
                UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
                {
                    ShowStr = $"当前持有的{GameUtil.Instance.GetItemName(itemId)}数量不足，是否前往获取？",
                    YesFun = () =>
                    {
                        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                        UITool.Open<UIGiftBag>(null, ActivityType.Element);
                    },
                });
                return;
            }

            isCanOption = false;
            OnJinJie();
        });

        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIFuNeng>();
        });
    }

    public override void OnEnter(params object[] args)
    {
        //JEvent.defaultEvent.Register(this);
        RegisterEvent();

        isAutoFuNeng = false;
        isCanOption = true;
        GetView().spineJinJie.color = Color.clear;

        GameUtil.Instance.SetItemIcon((int)CFG_Item.ItemID.ElementFN, GetView().img_FuNeng);
        GameUtil.Instance.SetItemIcon((int)CFG_Item.ItemID.ElementTP, GetView().img_JinJie);
        GameUtil.Instance.SetItemIcon((int)CFG_Item.ItemID.Gold, GetView().img_Gold);

        selectedElementUIPos = GetDefalutSelectedElementUIPos();
        ShowElementsInfo(true);
        SelectElement(selectedElementUIPos);
    }

    private Vector3[] attributeShowPosition = 
    {
        new Vector3(-122, 155), // uipos1
        new Vector3(200, 14),   // uipos2
        new Vector3(186, -68),  // uipos3
        new Vector3(-185, -70), // uipos4
        new Vector3(-200, 8),   // uipos5
        new Vector3(-116, 275), // uipos6
        new Vector3(250, 26),   // uipos7
        new Vector3(268, -185), // uipos8
        new Vector3(-266, -185),// uipos9
        new Vector3(-238, 26),  // uipos10
    };

    // 从火强开始数顺时针到外圈结束的元素位置顺序。自动赋能用。
    // 元素UI位置转为槽位id
    private int[] elementsUIPosToId = { 4, 1, 2, 5, 3, 9, 6, 7, 10, 8 };
    // 元素槽位id转为UI位置
    //private int[] elementsIdToUIPos = { 2, 3, 5, 1, 4, 7, 8, 10, 6, 9 };

    // 选中的元素ui中的位置
    private int selectedElementUIPos;
    // 当前选择的元素是否为进阶阶段
    private bool isNeedJinJie;
    // 是否是自动赋能
    private bool isAutoFuNeng;
    // 选中的元素当前等级配置
    private CFGItem_Element cfgCurrentElement;
    // 选中的元素下一等级配置
    private CFGItem_Element cfgNextElement;
    // 选中的元素的最大等级
    private int maxLevelCurrentElement;
    // 选中的元素的当前等级
    private int currentElementLevel;
    // 自动赋能间隔时间
    const float autoWaitTime = 0.35f;
    // 赋能结果返回前无法操作
    private bool isCanOption = true;

    private void ShowRes()
    {
        GetView().root_resGold.gameObject.SetActive(!isNeedJinJie);
        GetView().root_resJinJie.gameObject.SetActive(isNeedJinJie);
        GetView().txt_FuNeng.text = GameUtil.Instance.FormatNumWithUnit(BagM.Instance.GetNumByCfgID((int)CFG_Item.ItemID.ElementFN));
        GetView().txt_Gold.text = GameUtil.Instance.FormatNumWithUnit(BagM.Instance.GetNumByCfgID((int)CFG_Item.ItemID.Gold));
        GetView().txt_JinJie.text = GameUtil.Instance.FormatNumWithUnit(BagM.Instance.GetNumByCfgID((int)CFG_Item.ItemID.ElementTP));
    }

    private void RegisterEvent()
    {
        UniGameEvent.AddEventListener<GameEvent.WuXingFuNengUpdate>(OnSync_FuNeng);
        UniGameEvent.AddEventListener<Event_Source_Refresh>(OnEvent_Source_Refresh);
    }

    private void UnregisterEvent()
    {
        UniGameEvent.RemoveEventListener<GameEvent.WuXingFuNengUpdate>(OnSync_FuNeng);
        UniGameEvent.RemoveEventListener<Event_Source_Refresh>(OnEvent_Source_Refresh);
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnEvent_Source_Refresh(Event_Source_Refresh _e)
    {
        if (isNullGo)
            return;

        ShowRes();
        ShowConsumeAndButton();
    }

    private void OnSync_FuNeng(GameEvent.WuXingFuNengUpdate _e)
    {
        if (isNullGo)
            return;

        if (_e.isBreak)
            PlayJinJieAnimation();

        ShowElementsInfo();

        // 如果正在自动赋能，只刷新等级和属性显示，其他选中信息延迟到下次自动选时刷新。
        if (isAutoFuNeng)
        {
            ShowAttributeInfo();
            ClearAutoFuNengItem();
            autoFuNengTimeItem = GameHelper.Instance.AddInvoke(new TimeItem(autoWaitTime, AotuFuNeng));
        }
        else
        {
            isCanOption = true;
            SelectElement(selectedElementUIPos);
        }
    }

    Tweener tweener;

    private void PlayJinJieAnimation()
    {
        ClearTween();
        GetView().spineJinJie.color = Color.white;
        SpineUtil.Instance.PlayOnce(GetView().spineJinJie, "animation", null);
        tweener = GetView().spineJinJie.DOFade(0, 0.3f).SetDelay(2);
    }

    private void ClearTween()
    {
        if (tweener != null)
        {
            tweener.Kill();
            tweener = null;
        }
    }

    private void AotuFuNeng()
    {
        selectedElementUIPos++;
        if (selectedElementUIPos > 10)
            selectedElementUIPos = 1;

        SelectElement(selectedElementUIPos);

        int slotId = elementsUIPosToId[selectedElementUIPos - 1];
        if (IsMaxLevel(slotId))
        {
            GetView().txt_AutoFuNeng.text = "一键赋能";
            isAutoFuNeng = false;
            UITool.ShowTips("元素已达到最大等级！");
            return;
        }

        if (isNeedJinJie)
        {
            GetView().txt_AutoFuNeng.text = "一键赋能";
            isAutoFuNeng = false;
            UITool.ShowTips("请进阶后再赋能！");
            return;
        }

        if (!IsEnoughConsumeFN(out int itemId))
        {
            GetView().txt_AutoFuNeng.text = "一键赋能";
            isAutoFuNeng = false;
            GetView().btn_FuNeng.gameObject.SetActive(true);

            UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
            {
                ShowStr = $"当前持有的{GameUtil.Instance.GetItemName(itemId)}数量不足，是否前往获取？",
                YesFun = () =>
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    if ((CFG_Item.ItemID)itemId == CFG_Item.ItemID.Gold)
                        UITool.Open<UIShop>(null, CFG_Store.StoreType.Coupon);
                    else
                        UITool.Open<UIGiftBag>(null, ActivityType.Element);
                },
            });

            return;
        }

        OnFuNeng();
    }

    TimeItem autoFuNengTimeItem;
    private void ClearAutoFuNengItem()
    {
        if(GameHelper.Instance != null && autoFuNengTimeItem != null)
        {
            GameHelper.Instance.RemoveInvoke(autoFuNengTimeItem);
            autoFuNengTimeItem = null;
        }
    }

    public bool IsMaxLevel(int posId)
    {
        return GetLevel(posId) >= CFG_Element.MaxLevelOfPart(posId);
    }

    private int GetLevel(int posId) 
    {
        return TurnBasedM.Instance.empowers[posId].Level;
    }

    /// <summary>
    /// 打开界面时默认选择的元素。
    /// </summary>
    /// <returns>槽位id</returns>
    private int GetDefalutSelectedElementUIPos()
    {
        int posId = elementsUIPosToId[0];
        int lv = GetLevel(posId);
        int uiPos = 1;

        for (int i = 1; i < elementsUIPosToId.Length; i++) 
        {
            posId = elementsUIPosToId[i];
            if (GetLevel(posId) > lv)
            {
                lv = GetLevel(posId);
                uiPos = i + 1;
            }
        }

        return uiPos;
    }

    private void ShowElementsInfo(bool isInit = false)
    {
        for (int i = 0; i < GetView().elements.Length; i++)
        {
            int posId = elementsUIPosToId[i];
            GetView().elements[i].txt_Lv.text = "Lv." + GetLevel(posId);
            if (isInit) GetView().elements[i].posId = posId;
            GetView().elements[i].main.SetActive(TurnBasedM.Instance.elementMain == posId);
        }
    }

    private UIViewFuNeng.Element GetElementItemByPosId(int posId)
    {
        return Array.Find(GetView().elements, v => v.posId == posId);
    }

    /// <summary>
    /// 选择元素
    /// </summary>
    /// <param name="uiPos">1 - 10 </param>
    private void SelectElement(int uiPos)
    {
        selectedElementUIPos = uiPos;

        for (int i = 0; i < GetView().elements.Length; i++)
            GetView().elements[i].Select(selectedElementUIPos == i + 1);

        int slotId = elementsUIPosToId[selectedElementUIPos - 1];
        currentElementLevel = GetLevel(slotId);

        if (currentElementLevel > 0)
            cfgCurrentElement = CFG_Element.GetByPosAndLevel(slotId, currentElementLevel);
        else
            cfgCurrentElement = null;

        maxLevelCurrentElement = CFG_Element.MaxLevelOfPart(slotId);

        if (currentElementLevel < maxLevelCurrentElement)
            cfgNextElement = CFG_Element.GetByPosAndLevel(slotId, currentElementLevel + 1);
        else
            cfgNextElement = null;

        isNeedJinJie = currentElementLevel > 0 && TurnBasedM.Instance.empowers[slotId].BreakLevel < cfgCurrentElement.Breaklevel;

        ShowAttributeInfo();
        ShowRes();
        ShowConsumeAndButton();
    }

    private void ShowAttributeInfo()
    {
        // 无属性时显示 无属性
        GetView().root_AttributeWindow.anchoredPosition3D = attributeShowPosition[selectedElementUIPos - 1];
        GetView().txt_NoAttribute.gameObject.SetActive(currentElementLevel == 0);
        GetView().txt_AttributeInfo.gameObject.SetActive(currentElementLevel != 0);
        if (currentElementLevel != 0)
        {
            StringBuilder sb = new StringBuilder();
            for (int i = 0; i < cfgCurrentElement.Levelvalue.Length; i++)
                sb.AppendLine(GetAttributeStr(cfgCurrentElement.Levelvalue[i].key, cfgCurrentElement.Levelvalue[i].value));
            GetView().txt_AttributeInfo.text = sb.ToString().TrimEnd('\n');
        }
    }

    private void ShowConsumeAndButton()
    {
        GetView().root_ConsumeLayout.SetActive(currentElementLevel < maxLevelCurrentElement);
        GetView().txt_MaxLevel.gameObject.SetActive(currentElementLevel == maxLevelCurrentElement);

        if (currentElementLevel < maxLevelCurrentElement && !isNeedJinJie)
            GetView().txt_Probalility.text = $"赋能成功率：{cfgNextElement.Rate / 100}%";

        GetView().txt_Probalility.gameObject.SetActive(currentElementLevel < maxLevelCurrentElement && !isNeedJinJie);

        GetView().root_Buttons.SetActive(currentElementLevel < maxLevelCurrentElement);
        if (currentElementLevel < maxLevelCurrentElement)
        {
            GetView().txt_ConsumeDesc.text = isNeedJinJie ? "进阶消耗：" : "赋能消耗：";
            GetView().img_ConsumeFuNeng2.gameObject.SetActive(!isNeedJinJie);
            GetView().txt_ConsumeFuNeng2.gameObject.SetActive(!isNeedJinJie);
            int idConsume1 = isNeedJinJie ? cfgCurrentElement.Breakconsume[0].id : cfgNextElement.Levelconsume[0].id;
            int numConsume1 = isNeedJinJie ? cfgCurrentElement.Breakconsume[0].num : cfgNextElement.Levelconsume[0].num;
            GameUtil.Instance.SetItemIcon(idConsume1, GetView().img_ConsumeFuNeng1);
            GetView().txt_ConsumeFuNeng1.text = GameUtil.Instance.SetAndCheckStr(GameUtil.Instance.FormatNumWithUnit(numConsume1),
                idConsume1,
                numConsume1);

            if (!isNeedJinJie)
            {
                GameUtil.Instance.SetItemIcon(cfgNextElement.Levelconsume[1].id, GetView().img_ConsumeFuNeng2);
                GetView().txt_ConsumeFuNeng2.text = GameUtil.Instance.SetAndCheckStr(GameUtil.Instance.FormatNumWithUnit(cfgNextElement.Levelconsume[1].num),
                    cfgNextElement.Levelconsume[1].id,
                    cfgNextElement.Levelconsume[1].num);
            }

            GetView().btn_FuNeng.gameObject.SetActive(!isNeedJinJie && !isAutoFuNeng);
            GetView().btn_AutoFuNeng.gameObject.SetActive(!isNeedJinJie);
            GetView().btn_JinJie.gameObject.SetActive(isNeedJinJie);
        }
    }

    private string GetAttributeStr(int attrId, int value)
    {
        var item = CFG_XMLConfig.xmlConfig.RolePanelList.Find(v => v.SID == attrId);
        return $"{item.Name}：{CFG_XMLConfig.GetAttrValueToStr(item.Type, value)}";
    }

    /// <summary>
    /// 赋能消耗的资源是否够
    /// </summary>
    /// <param name="notEnoughItemId"></param>
    /// <returns></returns>
    private bool IsEnoughConsumeFN(out int notEnoughItemId)
    {
        if(!GameUtil.Instance.ComCoinToCheck(cfgNextElement.Levelconsume[0].id, cfgNextElement.Levelconsume[0].num, null, false))
        {
            notEnoughItemId = cfgNextElement.Levelconsume[0].id;
            return false;
        }

        if(!GameUtil.Instance.ComCoinToCheck(cfgNextElement.Levelconsume[1].id, cfgNextElement.Levelconsume[1].num, null, false))
        {
            notEnoughItemId = cfgNextElement.Levelconsume[1].id;
            return false;
        }

        notEnoughItemId = 0;
        return true;
    }

    /// <summary>
    /// 进阶消耗的资源是否够
    /// </summary>
    /// <param name="notEnoughItemId"></param>
    /// <returns></returns>
    private bool IsEnoughConsumeJinJie(out int notEnoughItemId)
    {
        if (!GameUtil.Instance.ComCoinToCheck(cfgCurrentElement.Breakconsume[0].id, cfgCurrentElement.Breakconsume[0].num, null, false))
        {
            notEnoughItemId = cfgCurrentElement.Breakconsume[0].id;
            return false;
        }
        notEnoughItemId = 0;
        return true;
    }

    /// <summary>
    /// 赋能
    /// </summary>
    private void OnFuNeng()
    {
        TcpC.C2S_ElementUpgradeRequest(elementsUIPosToId[selectedElementUIPos - 1]);
    }

    /// <summary>
    /// 进阶
    /// </summary>
    private void OnJinJie()
    {
        TcpC.C2S_ElementBreakthroughRequest(elementsUIPosToId[selectedElementUIPos - 1]);
    }

    public override void OnExit()
    {
        //JEvent.defaultEvent.Unregister(this);
        UnregisterEvent();
        ClearTween();
        ClearAutoFuNengItem();
        isAutoFuNeng = false;
        isCanOption = true;
        GetView().txt_AutoFuNeng.text = "一键赋能";
    }

    public override void OnDestroy()
    {
        ClearAutoFuNengItem();
    }
}