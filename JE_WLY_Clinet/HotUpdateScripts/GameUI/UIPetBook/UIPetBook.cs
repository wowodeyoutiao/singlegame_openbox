using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using JEngine.UI;
using SuperScrollView;
using HotUpdateScripts;

public class UIPetBook : UIBase<UIModelPetBook, UIViewPetBook>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }
    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIPetBook>();
        });
        GetView().dropdown_Qulite.AddOptions(GetModel().DroStrList);
        GetView().dropdown_Qulite.onValueChanged.AddListener(SetSortQulite);

        GetView().infiniteScroll.InitGridView(0, OnFillItem);

    }

    public override void OnEnter(params object[] args)
    {
        select_Qulite = 0;
        GetView().dropdown_Qulite.value = 0;
        SetSortQulite(select_Qulite);
    }

    private LoopGridViewItem OnFillItem(LoopGridView gridView, int sordIdx, int row, int column)
    {
        LoopGridViewItem _LoopGridViewItem = gridView.NewListViewItem("Item");
        GameObject go = _LoopGridViewItem.gameObject;

        JUIExtensions.btn(go.transform).onClick.RemoveAllListeners();
        var _item = siftList[sordIdx];
        var petInfo = PetM.Instance.GetOnePetByCfgId(_item.ID);
        go.transform.Find("unGet").gameObject.SetActive(!PetM.Instance.isHavePet(_item.ID));
        JUIExtensions.btn(go.transform).onClick.AddListener(() => {
            UITool.Open<UIPetShow>(null, _item.ID, petInfo == null ? 1 : petInfo.Level, petInfo == null ? 0 : petInfo.Star, petInfo != null);
        });
        GameUtil.Instance.SetItemIcon(_item.ItemID, JUIExtensions.img(go.transform, "icon"), JUIExtensions.img(go.transform), null,_item.Quality);

        JUIExtensions.img(go.transform, "img_pet").enabled = _item.Mark != 0;
        ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.CommonAtlas, "pet_" + _item.Mark, JUIExtensions.img(go.transform, "img_pet"));

        return _LoopGridViewItem;
           
    }

    private int select_Qulite = 0;
    private void SetSortQulite(int idx)
    {
        select_Qulite = GetModel().DroQualityList[idx];
        SetSiftList();
    }

    private List<CFGItem_Pet> siftList = new List<CFGItem_Pet>();
    private void SetSiftList()
    {
        siftList.Clear();
       
        for (int i = 0; i < CFG_Pet.Get().DicList.Count; i++)
        {
            if (CFG_Pet.Get().DicList[i].Quality >= select_Qulite || select_Qulite == 0)
            {
                siftList.Add(CFG_Pet.Get().DicList[i]);
            }
        }
        siftList.Sort((a, b) => {
            int qualityComparison = b.Quality.CompareTo(a.Quality);
            if (qualityComparison != 0)
                return qualityComparison;

            bool aExists = PetM.Instance.GetOnePetByCfgId(a.ID) != null;
            bool bExists = PetM.Instance.GetOnePetByCfgId(b.ID) != null;

            if (aExists != bExists)
                return aExists ? -1 : 1;

            return b.ID.CompareTo(a.ID);

        });

        if (siftList.Count != GetView().infiniteScroll.ItemTotalCount)
            GetView().infiniteScroll.SetListItemCount(siftList.Count);
        GetView().infiniteScroll.RefreshAllShownItem();
    }

    public override void OnExit()
    {
        
    }

     public override void OnDestroy()
    {

    }
}