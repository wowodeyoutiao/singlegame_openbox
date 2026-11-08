using JEngine.UI;

using MetaFramework.UI;
using UnityEngine;
using UnityEngine.UI;


public class HallEquipTabBtnItem
{
    public Button btn;
    public GameObject selGo;
    public Text text;

    public HallEquipTabBtnItem(Transform trans)
    {
        btn = trans.GetComponent<Button>();
        selGo = trans.Find("SelImg").gameObject;
        text = trans.Find("Text").GetComponent<Text>();
    }
}

public class HallEequipItem
{
    public Button btn;
    public Image bgImg;
    public GameObject silhouetteGo; // 剪影
    public Image iconImg;
    public Text equipLvText;
    public GameObject lockGo;
    public Text lvText;
    public GameObject starGo;
    public Text starText;
    public Image[] gemstoneImgs;
    public Transform baseEffectTrans;
    public GameObject baseEffectGo;

    public HallEequipItem(Transform trans)
    {
        btn = trans.GetComponent<Button>();
        iconImg = trans.Find("icon").GetComponent<Image>();
        bgImg = trans.GetComponent<Image>();
        silhouetteGo = trans.Find("bgImg").gameObject;
        lvText = trans.Find("txt_lv").GetComponent<Text>();
        starGo = trans.Find("StarIcon").gameObject;
        starText = trans.Find("StarIcon/txt_Starlv").GetComponent<Text>();
        equipLvText = trans.Find("icon/txt_lv").GetComponent<Text>();

        var stonsTrans = trans.Find("inset");
        gemstoneImgs = new Image[stonsTrans.childCount];
        for(int i = 0; i < stonsTrans.childCount;i++)
            gemstoneImgs[i] = stonsTrans.GetChild(i).GetComponent<Image>();

        baseEffectTrans = trans.Find("BestEffect");
        baseEffectGo = baseEffectTrans.gameObject;
    }
}

public class HallLevelItem
{
    public Button btn;
    public GameObject finishGo;
    public GameObject monsterNorIconGo;
    public GameObject monsterEliteIconGo;
    public GameObject monsterPassBossIconGo;
    public GameObject monsterBossIconGo;
    public GameObject monsterBossIconGo2;
    public GameObject attackGo;
    public GameObject challengeGo;

    public HallLevelItem(Transform trans)
    {
        btn = trans.GetComponent<Button>();
        finishGo = trans.Find("FinishIcon").gameObject;
        monsterNorIconGo = trans.Find("MonsterIcon1").gameObject;
        monsterEliteIconGo = trans.Find("MonsterIcon2").gameObject;
        monsterPassBossIconGo = trans.Find("MonsterIcon3").gameObject;
        monsterBossIconGo = trans.Find("MonsterIcon4").gameObject;//已解锁的
        monsterBossIconGo2 = trans.Find("MonsterIcon5").gameObject;// 未解锁的
        attackGo = trans.Find("AttackIcon").gameObject;
        var challengeTrans = trans.Find("ChallengeNode");
        if (challengeTrans != null)
            challengeGo = challengeTrans.gameObject;
    }
}

public class UIViewHall : MonoBehaviour, IUIView
{
	//背景
	public Transform bgFarNodeTrans;
    public Transform bgMidNodeTrans;
    public Transform bgNearbyNodeTrans;

    public GameObject restBgGo;

    // 角色与怪
    public Transform rolesNodeTrans;
    [HideInInspector]
    public Transform[] roleTranss;
    public Transform monsterTrans;

    public Transform roleLvUpTrans;
    public Text roleLvUpText;

    public Transform normalRewardTrans;
    public RectTransform normalRewardRect;

    // 宠物
    public Transform petTrans;
    public Transform petTalkTrans;
    [HideInInspector]
    public RectTransform petTalkRectTrans;
    public Text petTalkText;
    public Image petTalkImg;

    // hp change
    public Transform hpChangeNodeTrans;

    public Transform topLeftBtnsTrans;

    // 右上角按钮
    public Transform topRightBtnsTrans;

    // 装备
    public Transform equipBtnsNodeTrans;
    [HideInInspector]
    public HallEquipTabBtnItem[] equipTabBtnItems;

    public Transform equipItemsNodeTrans;
    [HideInInspector]
    public HallEequipItem[] equipItems;

    // 宝箱
    public GameObject boxShowUI;
    public Transform multipleBoxOpenTrans;
    public GameObject oneBoxOpenGo;
    public Button boxDelBtn;
    public Button boxClothingBtn;

    // 任务
    public GameObject taskNormalTaskInfoGo;
    public GameObject taskEndTaskInfoGo;
    public Button taskBtn;
    public Text taskProText;
    public Text taskDesText;
    public GameObject taskStateGo;
    //public Transform tfTaskFinishEffect;
    public Image taskImg;
    public Image taskBgImg;
    public GameObject taskGuideGo;

    // 体力
    public Button mainPowerBtn;
    public Slider mainPowerSlider;
    public Text mainPowerSliderText;
    public Text mainPowerText;

    // 开箱设置
    public Button boxOpenSetBtn;
    public Transform rootBoxSet;
    public Transform rootBoxSetOptions;
    public Button boxOpenSetSaveBtn;
    public Button boxOpenSetCloseBtn;
    public Button boxOpenSetMaskBtn;
    [HideInInspector]
    public Toggle[] boxSetToggles;

    // 聊天
    public Transform chatTrans;
    public Button chatBtn;

    // 关卡
    public Text levelTitleText;
    public Image levelSliderImg;
    public Transform levelRestItem;
    public Transform levelNodeTrans;
    [HideInInspector]
    public HallLevelItem[] levelItems;

    // loading
    public Transform BattleToRestLoadingTrans;
    public Transform ZhuanLoadingTrans;

    public Transform battleBossRewardTrans;

    // 普通血条
    public GameObject norMonsterHPHudGo;
    public Slider norMonsterHPSlider;

    // boss 血条
    public GameObject boosMonsterHPHudGo;
    public Slider bossMonsterHPSlider1;
    public Slider bossMonsterHPSlider2;
    public Image bossMonsterHPImg1;
    public Image bossMonsterHPImg2;
    public Text bossMonsterScheduleText;
    public Text bossMonsterNumText;
    public Text bossMonsterTitleText;
    public Button bossMonsterHPBtn;

    // 怪物大招
    public Transform bossMasterStrokeTrans;

    // boss 来袭
    public Transform bossCommingTrans;

    public Button guideBoxBtn;

    public Transform summaryTrans;
    // 右上角伸缩
    [HideInInspector]
    public Transform summaryBtnTrans;
    [HideInInspector]
    public Transform sevenDayBtnTrans; // 七天
    [HideInInspector]
    public Transform firstPayBtnTrans; // 首充
    [HideInInspector]
    public Transform battleSupportGiftTrans; // 助战礼包
    [HideInInspector]
    public Transform limitetFightRankTrans; // 限时排行（战役排行）
    [HideInInspector]
    public Transform limitetMainRankTrans; // 限时排行（主线排行）

    // 按钮
    [HideInInspector]
    public Transform specialOfferBtnTrans; // 特惠
    [HideInInspector]
    public Transform giftBagBtnTrans; // 礼包
    [HideInInspector]
    public Transform dailyBtnTrans; // 日常
    [HideInInspector]
    public Transform rankBtnTrans; // 排行

    // 连击
    public Transform comboTrans;
    [HideInInspector]
    public GameObject comboGo;
    public Text comboText;
    public GameObject comboEffectGo;

    public Transform effectRoot;

    // 设置
    public GameObject boxSettingBtnGo;

    // 控制按钮
    public Button backHomeBtn;
    public Button boxBtn;
    public Button attackBtn;

    // 大招
    public Transform skillBtnsTrans;
    public Transform skillShowTrans;

    // 特技
    public GameObject specialSkillNameGo;

    // 开箱子
    public Transform boxNodeTrans;
    [HideInInspector]
    public Button boxSpineBtn;
    [HideInInspector]
    public Text boxNumText;

    // 开箱子遮罩
    public GameObject boxShowUIMaskGo;

    public GameObject goBlessRedDotRoot;
    public GameObject goBlessFuncNewFlag;
    public GameObject goBlessRedDot;

    // 自动释放特技
    public GameObject autoSpecialSkillGo;
    public Button autoSpecialSkillBtn;

    public void Awake()
    {
        // 装备
        equipTabBtnItems = new HallEquipTabBtnItem[equipBtnsNodeTrans.childCount];
        for (int i = 0; i < equipBtnsNodeTrans.childCount; i++)
            equipTabBtnItems[i] = new HallEquipTabBtnItem(equipBtnsNodeTrans.GetChild(i));

        equipItems = new HallEequipItem[equipItemsNodeTrans.childCount];
        for (int i = 0; i < equipItemsNodeTrans.childCount; i++)
            equipItems[i] = new HallEequipItem(equipItemsNodeTrans.GetChild(i));

        // 角色
        roleTranss = new Transform[rolesNodeTrans.childCount];
        for (int i = 0; i < rolesNodeTrans.childCount; i++)
        {
            roleTranss[i] = rolesNodeTrans.Find("pos_" + i);
        }

        boxSetToggles = new Toggle[6];
        for (int i = 0; i < 6; i++)
        {
            boxSetToggles[i] = JUIExtensions.toggle(rootBoxSetOptions.GetChild(i));
        }

        petTalkRectTrans = petTalkTrans.GetComponent<RectTransform>();

        // 关卡
        levelItems = new HallLevelItem[5];
        for (int i = 0; i < 5; i++)
            levelItems[i] = new HallLevelItem(levelNodeTrans.Find($"LevelItem{i + 1}"));

        normalRewardRect = normalRewardTrans.GetComponent<RectTransform>();

        summaryBtnTrans = topRightBtnsTrans.Find("Btn_1");
        sevenDayBtnTrans = topRightBtnsTrans.Find("Btn_2");
        firstPayBtnTrans = topRightBtnsTrans.Find("Btn_3");
        limitetMainRankTrans = topRightBtnsTrans.Find("Btn_4");
        limitetFightRankTrans = topRightBtnsTrans.Find("Btn_5");
        battleSupportGiftTrans = topRightBtnsTrans.Find("Btn_6");

        dailyBtnTrans = summaryTrans.Find("Btns/Btn_1");
        rankBtnTrans = summaryTrans.Find("Btns/Btn_2");
        specialOfferBtnTrans = summaryTrans.Find("Btns/Btn_3");
        giftBagBtnTrans = summaryTrans.Find("Btns/Btn_4");

        comboGo = comboTrans.gameObject;

        boxSpineBtn = boxNodeTrans.Find("BoxSpine").GetComponent<Button>();
        boxNumText = boxNodeTrans.Find("NumBgImg/NumText").GetComponent<Text>();
        
        goBlessRedDotRoot = boxBtn.transform.Find("redDotRoot").gameObject;
        goBlessFuncNewFlag = boxBtn.transform.Find("functionOpenNewFlag").gameObject;
        goBlessRedDot = boxBtn.transform.Find("redDotRoot/RedImg").gameObject;

        goBlessFuncNewFlag?.SetActive(false);
        goBlessRedDotRoot?.SetActive(true);
    }
    
	public void Init(GameObject go)
	{

	}

    public void OnDestroy()
    {
		UITool.RemoveUI<UIHall>();
	}
}