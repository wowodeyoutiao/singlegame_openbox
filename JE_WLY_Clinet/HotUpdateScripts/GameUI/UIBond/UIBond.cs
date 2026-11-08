using BestHTTP.SecureProtocol.Org.BouncyCastle.Asn1.Ocsp;

using HotUpdateScripts.Config;
using HotUpdateScripts.ProtoM;

using JEngine.Event;
using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;

using pb;

using SuperScrollView;

using System;
using System.Collections.Generic;
using System.Linq;

using UnityEngine;
using UnityEngine.UI;

public class UIBond : UIBase<UIModelBond, UIViewBond>
{
    private UIViewBond uIViewBond;
    private Dictionary<int, BondItem> partnerLinkDict = new Dictionary<int, BondItem>();
    private Dictionary<GameObject, int> itemSelGoDict = new Dictionary<GameObject, int>();

    private List<CFGItem_PartnerLink> partnerLinkInfoList = new List<CFGItem_PartnerLink>();

    private Color activeColor = new Color(77f / 255, 174f / 255, 31f / 255);
    private Color lockColor = new Color(150f / 255, 150f / 255, 150f / 255);

    private int curSelIndex = -1;
    private BondItem curShowBondItem;

    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        uIViewBond = GetView();

        uIViewBond.closeBtn.onClick.AddListener(delegate
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIBond>();
        });

        uIViewBond.helpBtn.onClick.AddListener(delegate
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            UITool.Open<UIComRule>(null, 44);
        });

        OnInitLeftView();
        OnInitRightView();
    }

    #region Left View
    private void OnInitLeftView()
    {
        uIViewBond.leftLoopListView.InitListView(0, LeftViewFill);
    }

    private void ShowLeftView()
    {
        partnerLinkDict = BondM.Instance.GetPartnerLinkDict();

        uIViewBond.leftLoopListView.SetListItemCount(partnerLinkDict.Count);
        uIViewBond.leftLoopListView.RefreshAllShownItem();
    }

    private LoopListViewItem2 LeftViewFill(LoopListView2 listView, int idx)
    {
        if (!partnerLinkDict.ContainsKey(idx + 1))
            return null;

        if (!partnerLinkDict.TryGetValue(idx + 1, out var bondItem))
            return null;

        if (bondItem.cfgData == null || bondItem.cfgData.PartnerIDs == null || bondItem.cfgData.PartnerIDs.Length == 0)
            return null;

        var partnerIDs = bondItem.cfgData.PartnerIDs;

        var item = listView.NewListViewItem("Item");
        item.name = "item" + idx;

        var selBgTrans = item.transform.Find("SelbgImg");
        selBgTrans.gameObject.SetActive(false);

        if (!itemSelGoDict.ContainsKey(selBgTrans.gameObject))
            itemSelGoDict.Add(selBgTrans.gameObject, idx);
        else
            itemSelGoDict[selBgTrans.gameObject] = idx;
        item.transform.Find("LogoImg").gameObject.SetActive(partnerIDs.Length >= 4);
        item.transform.Find("NameTxt").GetComponent<Text>().text = bondItem.cfgData.BondName;

        if (partnerIDs.Length >= 4)
            item.GetComponent<RectTransform>().sizeDelta = new Vector2(312, 244);
        else
            item.GetComponent<RectTransform>().sizeDelta = new Vector2(312, 156);

        var newGo = item.transform.Find("NGo").gameObject;

        var rolesTrans = item.transform.Find("Roles");
        var unLock = true;
        for (int i = 0; i < rolesTrans.childCount; i++)
        {
            var roleItem = rolesTrans.GetChild(i);
            roleItem.gameObject.SetActive(i < partnerIDs.Length);

            if(i < partnerIDs.Length)
            {
                var roleInfo = PartnerM.Instance.GetOnePartnerInfoByCFGID(partnerIDs[i]);
                var roleCfg = CFG_Partner.GetById(partnerIDs[i]);

                var iconImg = JUIExtensions.img(roleItem.transform, "IconImg");
                var bgImg = JUIExtensions.img(roleItem.transform);

                var roleBtn = roleItem.GetComponent<Button>();
                roleBtn.onClick.RemoveAllListeners();
                GameUtil.Instance.SetItemIcon(roleCfg.ItemID, iconImg, bgImg, null, roleCfg.Quality);
                var starNum = 0;
                if(roleInfo != null)
                {
                    starNum = roleInfo.Star;

                    //roleItem.Find("NumText").GetComponent<Text>().text = roleInfo.Star.ToString();
                    roleBtn.onClick.AddListener(() => 
                    {
                        OnSelectLeftItemAtIndex(idx);
                    });

                    roleBtn.onClick.AddListener(() =>
                    {
                        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                        Action action = () =>
                        {
                            UITool.Close<UIItemShow>();
                            UITool.Open<UIPartnerUpStar>(null, roleInfo.Uid, true);
                        };
                        UITool.Open<UIItemShow>(null, new UIItemShow.Params
                        {
                            cfgId = roleCfg.ItemID,
                            isIgnoreJob = true,
                            callback = action,
                            buttonName = "升星"
                        });
                    });
                }
                else
                {
                    unLock = false;
                    //roleItem.Find("NumText").GetComponent<Text>().text = "";
                    roleBtn.onClick.AddListener(() =>
                    {
                        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                        newGo.SetActive(false);
                        BondM.Instance.RemoveNewRed(bondItem.bondID);

                        Action action = () =>
                        {
                            if (CFG_XMLConfig.GetOpenSolt(RoleM.Instance.Level) > 0)
                            {
                                HttpC.GetPrivilegeDetail(1, (_r) =>
                                {
                                    UITool.Close<UIItemShow>();
                                    UITool.Open<UIPartnerCall>(null, _r);
                                });
                            }
                            else
                                UITool.ShowTips($"{RoleM.Instance.Level}级开启！");
                        };

                        UITool.Open<UIItemShow>(null, new UIItemShow.Params
                        {
                            cfgId = roleCfg.ItemID,
                            isIgnoreJob = true,
                            callback = action,
                            buttonName = "获取"
                        });
                    });
                }

                roleItem.Find("Star").gameObject.SetActive(starNum > 0);
                if (starNum > 0)
                    roleItem.Find("NumText").GetComponent<Text>().text = starNum.ToString();
                else
                    roleItem.Find("NumText").GetComponent<Text>().text = "";

                iconImg.material = roleInfo != null ? null : uIViewBond.grayImg.material;
                bgImg.material = roleInfo != null ? null : uIViewBond.grayImg.material;
            }
        }

        newGo.SetActive(unLock && BondM.Instance.IsNew(bondItem.bondID));

        //Debug.LogError($"index:{idx}, unLock: {unLock}, isNew： {isNew}");

        var btn = item.GetComponent<Button>();
        btn.onClick.RemoveAllListeners();
        btn.onClick.AddListener(() =>
        {
            newGo.SetActive(false);
            BondM.Instance.RemoveNewRed(bondItem.bondID);
            OnSelectLeftItemAtIndex(idx, false);
        });

        return item;
    }

    private void OnSelectLeftItemAtIndex(int index, bool froce = false)
    {
        AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
        if (curSelIndex == index && !froce) return;
        curSelIndex = index;

        //Debug.LogError("OnSelectLeftItemAtIndex" + index);
        foreach (var item in itemSelGoDict)
        {
            //Debug.LogError($"{item.Key}, {item.Value}, {index}");
            item.Key.SetActive(item.Value == index);
        }

        ShowRightView(index);
    }
    #endregion

    #region right view
    private void OnInitRightView()
    {
        uIViewBond.rightLoopListView.InitListView(0, RightViewFill);
    }

    private void ShowRightView(int index)
    {
        curShowBondItem = null;
        if (partnerLinkDict.TryGetValue(index + 1, out var bondItem))
        {
            //Debug.LogError("ShowRightView, configId:  " + bondItem.cfgData.ID);
            uIViewBond.rightTitleText.text = bondItem.cfgData.BondName;
            curShowBondItem = bondItem;
            partnerLinkInfoList = CFG_PartnerLink.GetPartnerLinkByBondId(bondItem.bondID);
            partnerLinkInfoList.Sort((x, y) => { return x.ID.CompareTo(y.ID); });

            uIViewBond.rightLoopListView.SetListItemCount(partnerLinkInfoList.Count);
            uIViewBond.rightLoopListView.RefreshAllShownItem();
        }
    }

    private LoopListViewItem2 RightViewFill(LoopListView2 listView, int idx)
    {
        if (partnerLinkInfoList.Count < idx)
            return null;

        var partnerLink = partnerLinkInfoList[idx];

        var item = listView.NewListViewItem("Item");
        item.name = "item" + idx;

        //Debug.LogError($"{partnerLink.BondID}, {curShowBondItem.cfgData.ID}, {curShowBondItem.cfgData.Star}, {partnerLink.Star}");
        var unlock = curShowBondItem.cfgData.Star >= partnerLink.Star && curShowBondItem.bondState == BondState.Active;

        var titleTxt = item.transform.Find("Info/DesTxt").GetComponent<Text>();
        if (idx == 0)
        {
            titleTxt.text = $"获得{partnerLink.PartnerIDs.Length}种伙伴";
            titleTxt.color = unlock ? activeColor : lockColor;
        }
        else
        {
            titleTxt.text = $"全部升至 {partnerLink.Star} 星";
            titleTxt.color = unlock ? activeColor : lockColor;
        }

        var stateTxt = item.transform.Find("Info/StateTxt").GetComponent<Text>();
        if (unlock)
        {
            stateTxt.text = $"已激活";
            stateTxt.color = activeColor;
        }
        else
        {
            stateTxt.text = $"未激活";
            stateTxt.color = lockColor;
        }

        var buffTrans = item.transform.Find("Buff");
        for (int i = 0; i < buffTrans.childCount; i++)
            buffTrans.GetChild(i).gameObject.SetActive(false);

        var buffDesGo = item.transform.Find("Buff/DesTxt").gameObject;
        for(int i = 0; i < Math.Max(partnerLink.Attribute.Length, buffTrans.childCount); i++)
        {
            if(i < partnerLink.Attribute.Length)
            {
                var attribute = partnerLink.Attribute[i].ElementAt(0);
                var attributeCfg = CFG_XMLConfig.GetAttributeByID(int.Parse(attribute.Key));
                if (attributeCfg != null)
                {
                    GameObject desGo;
                    if (i < buffTrans.childCount)
                        desGo = buffTrans.GetChild(i).gameObject;
                    else 
                        desGo = GameObject.Instantiate(buffDesGo, buffDesGo.transform.parent);
                    desGo.SetActive(true);
                    var desText = desGo.GetComponent<Text>();
                    desText.text = $"{attributeCfg.Name} + {CFG_XMLConfig.GetAttrValueToStr(attributeCfg.Type, attribute.Value)}";
                    desText.color = unlock ? activeColor : lockColor;
                }
            }
        }
        item.GetComponent<RectTransform>().sizeDelta = new Vector2(0, 40 + 26 * partnerLink.Attribute.Length);
        return item;
    }
    #endregion

    public override void OnEnter(params object[] args)
    {
        //JEvent.defaultEvent.Register(this);
        UniGameEvent.AddEventListener<GameEvent.Refrsh_Common>(OnRefresh_CommonResponse);

        curSelIndex = -1;
        ShowLeftView();
        OnSelectLeftItemAtIndex(0);
    }

    //[Subscriber(ThreadMode.Main)]
    private void OnRefresh_CommonResponse(GameEvent.Refrsh_Common _e)
    {
        if (isNullGo && _e.cid == CmdId.C2S_PartnerStar && _e.code != 0) return;

        if (curSelIndex == -1) curSelIndex = 0;
        ShowLeftView();
        OnSelectLeftItemAtIndex(curSelIndex, true);
    }

    public override void OnExit()
    {
        //JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refrsh_Common>(OnRefresh_CommonResponse);
    }

    public override void OnDestroy()
    {
        //JEvent.defaultEvent.Unregister(this);
        UniGameEvent.RemoveEventListener<GameEvent.Refrsh_Common>(OnRefresh_CommonResponse);

    }
}