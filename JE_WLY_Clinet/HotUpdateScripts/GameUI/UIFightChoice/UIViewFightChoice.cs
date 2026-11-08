using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using SuperScrollView;

public class UIViewFightChoice : MonoBehaviour,IUIView
{
	public Text txt_title;
	public Button btn_help;
	public Transform root_SelectBtns;

    public Button btn_Close;

    public Transform root_ChallengeBtns;
    public Transform root_UI;
    public GameObject obj_ZY;
    public GameObject obj_YJ;
    public GameObject obj_SLTower;
    public GameObject obj_SL;
    public GameObject obj_PvP;

    public Text txt_ZYName;
	public Transform root_ZYAward;
	public Image img_ZYBoss;
    public Text txt_ZYlimet;
    public Button btn_ZYPlay;
    public Button btn_FightAward;

    public Transform root_YJLevel;
    public Button btn_YJ_Add;
    public Button btn_YJ_Jian;
    public Text txt_YJ_Name;
    public Text txt_YJ_Num;
    public Button btn_YJ_CountAdd;
    public Button btn_YJ_Play;
    public Button btn_YJ_SD;
    public Slider slider_YJ;
    public GameObject obj_YJLock;
    public Text txt_YJ_StarNum;


    public LoopListView2 InfiniteScroll_SLTower;
    public Button btn_SLTower_Paly;
    public Button btn_SLTower_SLShop;
    public Button btn_SLTower_SLJJ;

    public LoopListView2 InfiniteScroll_SL;
    public Text txt_SL_Count;
    public Button btn_SL_CountAdd;

    public Button btn_PVP_ZhanHun;
    public Transform root_PvPOne;
    public Transform root_PvPMore;

    public Transform Obj_HaiDao;
    

    public void Awake()
    {

    }
	public void Init(GameObject go)
	{

	}
	public void OnDestroy()
	{
		 UITool.RemoveUI<UIFightChoice>();
	}
}