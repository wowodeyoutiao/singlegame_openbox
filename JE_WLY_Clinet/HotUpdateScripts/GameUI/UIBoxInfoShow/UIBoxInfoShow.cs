using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;

using System.Collections;
using System.Collections.Generic;

using UnityEngine;
using UnityEngine.UI;

public class UIBoxInfoShow : UIBase<UIModelBoxInfoShow, UIViewBoxInfoShow>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIBoxInfoShow>();
        });

        for (int i = 0; i < 4; i++)
        {
            int idx = i;
            JUIExtensions.btn(GetView().awardRoot.Find("Item_" + i)).onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                UITool.Open<UIItemShow>(null, new UIItemShow.Params { cfgId = BoxInfo[selectItemListIndex][idx].Id });
            });
        }
    }

    public struct BoxItem {
        public int Id;
        public int Num;
    }

    private List<BoxItem[]> BoxInfo;
    private int selectItemListIndex = 0;

    public struct Buttons
    {
        public string Name;
    }

    public List<Buttons> buttons = new List<Buttons>();


    public override void OnEnter(params object[] args)
    {
        buttons = args[0] as List<Buttons>;
        SetButtons();
        BoxInfo = args[1] as List<BoxItem[]>;

        if (args.Length >= 3)
            GetView().txt_Tip.text = args[2] as string;
        else
            GetView().txt_Tip.text = "";

        SelectButton(0);
    }

    private void SetButtons()
    {
        bool hasButtons = buttons != null && buttons.Count > 0;
        GetView().buttonRoot.gameObject.SetActive(hasButtons);
        if (!hasButtons)
            return;

        for (int i = 0; i < GetView().buttonRoot.childCount; i++)
        {
            GetView().buttonRoot.GetChild(i).gameObject.SetActive(false);
        }

        for (int i = 0; i < buttons.Count; i++)
        {
            if(i >= GetView().buttonRoot.childCount)
            {
                Object.Instantiate(GetView().buttonRoot.GetChild(0).gameObject, GetView().buttonRoot);
            }
            var btn = GetView().buttonRoot.GetChild(i);
            JUIExtensions.txt(btn.Find("Label")).text = buttons[i].Name;
            JUIExtensions.txt(btn.Find("Select/Label")).text = buttons[i].Name;
            btn.gameObject.SetActive(true);
            int index = i;
            JUIExtensions.btn(btn).onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                SelectButton(index);
            });
        }
    }

    private void SelectButton(int index)
    {
        for (int i = 0; i < GetView().buttonRoot.childCount; i++)
        {
            var btn = GetView().buttonRoot.GetChild(i);
            btn.Find("Select").gameObject.SetActive(i == index);
        }
        selectItemListIndex = index;
        ShowItems();
    }

    private void ShowItems()
    {
        for (int i = 0; i < 4; i++)
        {
            var items = BoxInfo[selectItemListIndex];
            var root = GetView().awardRoot.Find("Item_" + i);
            root.gameObject.SetActive(i < items.Length);
            if (i < items.Length)
            {
                GameUtil.Instance.SetItemIcon(items[i].Id, JUIExtensions.img(root.Find("icon")), JUIExtensions.img(root, "bg"), JUIExtensions.img(root, "icon/nature"));
                JUIExtensions.txt(root.Find("txt_num")).text = items[i].Num.ToString();
            }
        }
    }

    public override void OnExit()
    {

    }

    public override void OnDestroy()
    {

    }
}