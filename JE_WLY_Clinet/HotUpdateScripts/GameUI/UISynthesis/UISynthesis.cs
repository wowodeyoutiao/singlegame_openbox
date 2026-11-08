using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using JEngine.UI;
using JEngine.Helper;
using HotUpdateScripts.ProtoM;
using JEngine.Event;

public class UISynthesis : UIBase<UIModelSynthesis, UIViewSynthesis>
{
    public class SelectSyn
    {
        public int id;
        public long Uid;
        public int num;
    }

    public static List<SelectSyn> selectSyns = new List<SelectSyn>();
    public enum UISelect 
    {
        Material,
        Jewel,
        Used,
    }

    private Dictionary<UISelect, string> TypeToStrDic = new Dictionary<UISelect, string>
    {
        [UISelect.Material] = "材料",
        [UISelect.Jewel] = "魂石",
        [UISelect.Used] = "道具",
    };

    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UISynthesis>();
        });

        GetView().btn_Help.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIComRule>(null, 16);
        });

        GetView().btn_AutoAdd.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnAutoAdd();
        });

        GetView().btn_Del.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            selectSyns.Clear();
            OnUISelectShow();
        });

        GetView().btn_Syn.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

            if (!GameUtil.Instance.ComCoinToCheck(usedId, usedNum, () => {
                UITool.Close<UISynthesis>();
                UITool.Close<UIBag>();
                UITool.Close<UIBuild_Equip>();
            }))
            {
                return;
            }
            TcpC.SetCompose(comList);

            string page = _uISelect switch
            {
                UISelect.Material => "材料合成界面",
                UISelect.Jewel => "魂石合成界面",
                UISelect.Used => "道具合成界面",
                _ => "",
            };
            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "首页",
                component = "背包",
                position = page,
                event_click = "一键合成",
            });
        });

        for (int i = 0; i < 3; i++)
        {
            int idx = i;
            JUIExtensions.btn(GetView().root_top.GetChild(i)).onClick.AddListener(() => {


                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);

                OnSelectUIShow((UISelect)idx);
            });
        }

        for (int i = 0; i < 3; i++)
        {
            int idx = i;
            JUIExtensions.btn(GetView().root_Center.Find("Item_" + i)).onClick.AddListener(() => {

                if (isAmin) return;

                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                if (selectSyns.Count <= idx)
                {
                    UITool.Open<UISynthesisSelect>(null, _uISelect, selectSyns);
                }
                else
                {
                    selectSyns.RemoveAt(idx);
                    OnUISelectShow();
                }
            });
        }
    }

    private UISelect _uISelect;
    private int SynCount = 3;

    public override void OnEnter(params object[] args)
    {
        JEvent.defaultEvent.Register(this);
        isAmin = false;
        selectSyns.Clear();
        if (args[0] != null)
        {
            var baseInfo = (pb.ItemInfo)args[0];
            selectSyns.Add(new SelectSyn { id = baseInfo.CfgId, Uid = baseInfo.Uid, num = 1 });
            if (CFG_Item.GetById(baseInfo.CfgId).SonType == 3)
                _uISelect = UISelect.Jewel;
            else if (CFG_Item.GetById(baseInfo.CfgId).SonType == 1)
                _uISelect = UISelect.Material;
            else
                _uISelect = UISelect.Used;
        }
        else
        {
            _uISelect = (UISelect)args[1];
        }
        OnSelectUIShow(_uISelect);
    }

    private void OnSelectUIShow(UISelect uISelect)
    {
        for (int i = 0; i < 3; i++)
        {
            GetView().root_top.GetChild(i).Find("on").gameObject.SetActive(uISelect == (UISelect)i);
        }
        if (_uISelect != uISelect)
        {
            selectSyns.Clear();
        }
        _uISelect = uISelect;

        OnUISelectShow();
    }

    public void OnUISelectShow()
    {
        for (int i = 0; i < 3; i++)
        {
            var go = GetView().root_Center.Find("Item_" + i);
            go.Find("icon").gameObject.SetActive(i < selectSyns.Count);
            go.Find("add").gameObject.SetActive(i >= selectSyns.Count);
            if (i < selectSyns.Count)
            {
                GameUtil.Instance.SetItemIcon(selectSyns[i].id, JUIExtensions.img(go, "icon"), JUIExtensions.img(go), JUIExtensions.img(go, "icon/nature"));
                JUIExtensions.txt(go, "icon/txt").text = GameUtil.Instance.FormatNumWithUnit(selectSyns[i].num);
            }
            else
            {
                GameUtil.Instance.SetItemBg(0, JUIExtensions.img(go));
            }
        }

        GetView().btn_AutoAdd.gameObject.SetActive(selectSyns.Count > 0 && selectSyns.Count < SynCount);


        GetView().txt_tip.text = selectSyns.Count == SynCount ? "" : (selectSyns.Count == 0 ? "请点击槽位选择需要合成的"+ TypeToStrDic[_uISelect] : "全部槽位放满后才可合成");
        GetView().obj_SelectEnd.SetActive(selectSyns.Count == SynCount);
        GetView().spine_anim.gameObject.SetActive(false);
        GetView().obj_lose.SetActive(false);
        GetView().obj_NorTip.SetActive(true);
        GetView().root_Center.Find("Item_3").gameObject.SetActive(false);

        if (selectSyns.Count == SynCount)
            OnSelectEnd();
    }

    private List<pb.ComposeInfo> comList = new List<pb.ComposeInfo>();

    private int usedNum = 0;
    private int usedId = 0;
    private void OnSelectEnd()
    {
        var compose = CFG_Compose.GetByRawId(selectSyns[0].id);
        GetView().txt_pro.text = GameUtil.Instance.GetPercentStr((float)compose.SuccessRate / 10000f);

        int num = 0;
        bool isSame = true;

        int cid = selectSyns[0].id;
        comList.Clear();
        for (int i = 0; i < selectSyns.Count; i++)
        {
            num += selectSyns[i].num;

            if (cid != selectSyns[i].id)
            {
                isSame = false;
            }
        }
        num /= 3;

        if (!isSame)
        {
            num = 1;
        }

        for (int i = 0; i < selectSyns.Count; i++)
        {
            comList.Add(new pb.ComposeInfo { Num = num, Uid = selectSyns[i].Uid });
        }

        usedId = compose.NeedGold.id;
        usedNum = compose.NeedGold.num * num;
        GetView().txt_money.text = GameUtil.Instance.SetAndCheckStr(GameUtil.Instance.FormatNumWithUnit(usedNum), usedId, usedNum);
    }


    private bool isAmin = false;
    [Subscriber(ThreadMode.Main)]
    private void OnAnimShow(GameEvent.Sync_ComposeResponse _e)
    {
        isAmin = true;
        GetView().spine_anim.gameObject.SetActive(true);
        SpineUtil.Instance.PlayOnce(GetView().spine_anim, "animation", () => { 
            GetView().spine_anim.gameObject.SetActive(false);
        });

        AudioManager.Instance.PlaySound(AudioKey.Sound.UI_ComposeEffect);
        GetView().obj_SelectEnd.SetActive(false);
        GameHelper.Instance.AddInvoke(new TimeItem(1.3f, () =>
        {
            if (!isNullGo)
                OnResult();
        }));
        void OnResult()
        {
            for (int i = 0; i < 3; i++)
            {
                var go = GetView().root_Center.Find("Item_" + i);
                go.Find("icon").gameObject.SetActive(false);
                go.Find("add").gameObject.SetActive(false);
            }
            selectSyns.Clear();
            bool isWin = _e.data.Code == 0;
            GetView().obj_lose.SetActive(!isWin);
            GetView().obj_NorTip.SetActive(!isWin);
            var winRoot = GetView().root_Center.Find("Item_3");
            winRoot.gameObject.SetActive(isWin);

            if (isWin)
            {
                GameUtil.Instance.SetItemIcon((int)_e.data.NewConpose[0].Uid, JUIExtensions.img(winRoot, "icon"), JUIExtensions.img(winRoot, "icon_bg"), JUIExtensions.img(winRoot, "icon/nature"));
                JUIExtensions.txt(winRoot, "icon/txt").text = GameUtil.Instance.FormatNumWithUnit(_e.data.NewConpose[0].Num);
                if (_e.data.NewConpose.Count > 1)
                {
                    UITool.Open<UISynthesisMoreAward>(null, _e.data.NewConpose);
                }
            }

            GameHelper.Instance.AddInvoke(new TimeItem(4, () =>
            {
                isAmin = false;
                if (!isNullGo)
                    OnUISelectShow();
            }));
        }
        
       
    }

    //自动添加
    public void OnAutoAdd(SelectSyn addSelect = null)
    {
        if (addSelect != null)
        {
            selectSyns.Clear();
            selectSyns.Add(addSelect);
        }
        BagM.Instance.GetAutoSelectEquipSyns(SynCount, ref selectSyns);
        OnUISelectShow();
    }

    public override void OnExit()
    {
        JEvent.defaultEvent.Unregister(this);
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}