using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using JEngine.UI;
using MetaFramework.UI;
using System;
using System.Collections.Generic;
using UnityEngine.UI;

public class UIAwards : UIBase<UIModelAwards, UIViewAwards>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btn_back.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            ClearJa();
            OnBack();
        });

        for (int i = 0; i < GetView().root_Awards.childCount; i++)
        {
            int idx = i;
            GetView().root_Awards.GetChild(i).GetComponent<Button>()?.onClick.AddListener(() => {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = _AwardsInfo.Awards[idx].id });
            });
        }
    }
    private void OnBack()
    {
        UITool.Close<UIAwards>();
        _AwardsInfo._CallBack?.Invoke();
    }

    private TimeItem Ja;
    private void ClearJa()
    {
        if (Ja != null)
        {
            GameHelper.Instance.RemoveInvoke(Ja);
            Ja = null;
        }

        for (int i = 0; i < TimeItemList.Count; i++)
        {
            GameHelper.Instance.RemoveInvoke(TimeItemList[i]);
        }
        TimeItemList.Clear();
    }

    public class AwardsInfo
    {
        public string des;
        public int TimeCount;
        public Action _CallBack;
        public AwardItem[] Awards;
        public class AwardItem { 
            public int id;
            public int num; 
            public int kind; }
    }

    private AwardsInfo _AwardsInfo;

    private Dictionary<int, string> KuangSpineName = new()
    {
        [2] = "lv",
        [3] ="lan",
        [4] ="zi",
        [6] = "cheng",
        [7] = "hong"
    };
    private List<TimeItem> TimeItemList = new List<TimeItem>();
    public override void OnEnter(params object[] args)
    {
        ClearJa();
        _AwardsInfo = args[0] as AwardsInfo;
        int showDelay = 0;
        for (int i = 0; i < GetView().root_Awards.childCount; i++)
        {
            var go = GetView().root_Awards.GetChild(i);
            go.gameObject.SetActive(false);
            showDelay = i ;
            if (i < _AwardsInfo.Awards.Length)
            {
                bool isSpineEff = KuangSpineName.ContainsKey((int)CFG_Item.GetById(_AwardsInfo.Awards[i].id).Quality) && CFG_Item.GetById(_AwardsInfo.Awards[i].id).Type == (int)CFG_Item.ItemType.PartnerCard;
                go.Find("spine").gameObject.SetActive(isSpineEff);
                if (isSpineEff)
                {
                    SpineUtil.Instance.PlayLoop(go.Find("spine"), KuangSpineName[(int)CFG_Item.GetById(_AwardsInfo.Awards[i].id).Quality]);
                }
                GameUtil.Instance.SetItemIcon(_AwardsInfo.Awards[i].id, JUIExtensions.img(go.transform.Find("icon")), JUIExtensions.img(go.transform), JUIExtensions.img(go.transform.Find("icon/nature")));
                JUIExtensions.txt(go, "txt_num").text = GameUtil.Instance.FormatNumWithUnit(_AwardsInfo.Awards[i].num);

                 var timer = GameHelper.Instance.AddInvoke(new TimeItem(showDelay * 0.5f,  () => {
                     if (isNullGo) return;
                        go.gameObject.SetActive(true);
                 }));
                TimeItemList.Add(timer);
            }
        }

        GetView().txt_des.text = _AwardsInfo.des;
        if (!GuidanceM.Instance.isGuiding)
        {
            int Count = _AwardsInfo.TimeCount;
            Ja = GameHelper.Instance.AddInvoke(new TimeItem(Count, 1, () =>
            {
                GetView().txt_Time.text = Count + "秒后将自动退出";
                Count--;
                if (Count <= 0)
                {
                    ClearJa();
                    OnBack();
                }
            }));
        }
        else GetView().txt_Time.text = "";
    }

    public override void OnExit()
    {
        ClearJa();
    }

     public override void OnDestroy()
    {
        OnExit();
    }
}