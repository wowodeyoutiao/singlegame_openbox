using Cysharp.Threading.Tasks;
using HotUpdateScripts;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using JEngine.UI;
using MetaFramework.UI;
using pb;
using Spine.Unity;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class UIPartnerBook : UIBase<UIModelPartnerBook, UIViewPartnerBook>
{
    private UIViewPartnerBook bookView;
    private string[] qualityNames = new string[] { "无双伙伴", "红色伙伴", "橙色伙伴", "紫色伙伴", "蓝色伙伴", "绿色伙伴" };
    private List<int> qualityList = new List<int>();
    private Dictionary<int, List<CFGItem_Partner>> partnerCfgDict = new Dictionary<int, List<CFGItem_Partner>>();
    private List<Transform> itemTransList = new List<Transform>();

    private Dictionary<int, Transform> roleTransDict = new Dictionary<int, Transform>();
    private Material originMat;

    private int selectItemId;

    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public override void OnInit()
    {
        bookView = GetView();

        bookView.closeBtn.onClick.AddListener(() => {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIPartnerBook>();
        });
    }

    public override void OnEnter(params object[] args)
    {
        bookView = GetView();

        var list = CFG_Partner.GetALlList();
        if(list == null || list.Count == 0)
        {
            UITool.Close<UIPartnerBook>();
            return;
        }
        selectItemId = 0;
        UniGameEvent.AddEventListener<GameEvent.Refrsh_Common>(OnRefresh_CommonResponse);
        bookView.contentGo.SetActive(false);
        qualityList.Clear();
        partnerCfgDict.Clear();

        for (int i = 0; i < list.Count; i++)
        {
            var cfgItem = list[i];
            if(!partnerCfgDict.ContainsKey(cfgItem.Quality))
            {
                partnerCfgDict.Add(cfgItem.Quality, new List<CFGItem_Partner>());
                qualityList.Add(cfgItem.Quality);
            }

            partnerCfgDict[cfgItem.Quality].Add(cfgItem);
        }

        qualityList.Sort();
        qualityList.Reverse();

        roleTransDict.Clear();
        for (int i = 0; i < qualityNames.Length; i++)
        {
            LoadSingleQuality(i, qualityList[i], qualityNames[i]);
        }
        bookView.contentGo.SetActive(true);
    }

    private void LoadSingleQuality(int index, int quality, string name)
    {
        if(partnerCfgDict.TryGetValue(quality, out var partners))
        {
            if (partners == null || partners.Count == 0) return;

            partners.Sort((a, b) => { return a.ID.CompareTo(b.ID); });

            Transform itemTrans;
            if(index < itemTransList.Count)
                itemTrans = itemTransList[index];
            else
            {
                itemTrans = bookView.GetItemTrans();
                itemTransList.Add(itemTrans);
            }

            itemTrans.localPosition = Vector3.zero;
            itemTrans.localScale = Vector3.one;

            var titleText = itemTrans.Find("TitleItem/TitleText").GetComponent<Text>();
            titleText.text = name;

            var roleContainerTrans = itemTrans.Find("RoleItem");
            if(roleContainerTrans.childCount > 1)
            {
                for(int i = roleContainerTrans.childCount - 1; i > 0; i--)
                {
                    var child = roleContainerTrans.GetChild(i);
                    GameHelper.Instance.ClearObjToPool(child);
                    GameObject.DestroyImmediate(child.gameObject);
                }
            }

            var go = roleContainerTrans.GetChild(0).gameObject;
            for (int i = 0; i < partners.Count; i++)
            {
                var partnersCfg = partners[i];
                var itemID = partnersCfg.ID;

                var canUp = false;

                var roleItemGo = GameObject.Instantiate(go);
                roleItemGo.transform.SetParent(roleContainerTrans);
                roleItemGo.transform.localScale = Vector3.one;
                roleItemGo.transform.localPosition = Vector3.zero;

                var partnerInfo = PartnerM.Instance.GetOnePartnerInfoByCFGID(itemID);
                var isHave = partnerInfo != null;

                if(isHave)
                    roleTransDict.Add(itemID, roleItemGo.transform);

                // 背景
                var bgImg = JUIExtensions.img(roleItemGo.transform.Find("BgImg"));
                if (isHave)
                    ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.PartnerCall, "bg_Summoning_dzk_" + partnersCfg.Quality, bgImg);
                else
                    ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.PartnerCall, "bg_Summoning_pzk", bgImg);

                var bgBtn = roleItemGo.transform.Find("BgImg").GetComponent<Button>();
                bgBtn.onClick.RemoveAllListeners();
                bgBtn.onClick.AddListener(() => 
                {
                    AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                    selectItemId = itemID;
                    if (canUp && partnerInfo != null)
                        UITool.Open<UIPartnerUpStar>(null, partnerInfo.Uid, true);
                    else
                        UITool.Open<UIStarDetails>(null, partnersCfg);
                });

                // 职业
                GameUtil.Instance.SetItemIcon(itemID, null, null, JUIExtensions.img(roleItemGo.transform.Find("NatureImg")));

                // 角色
                string poolKey = "Partner_" + itemID;
                var path = GameConst.GetMonsterModelByID(poolKey);
                GameHelper.Instance.SetAsyncObj(poolKey, path, (role) =>
                {
                    if (isNullGo)
                        return;

                    role.transform.SetParent(roleItemGo.transform.Find("IconImg"));
                    role.name = "SpineItem";
                    role.transform.localScale = Vector3.one * 0.32f;
                    role.transform.localPosition = new Vector3(0, -26, 0);
                    role.SetActive(true);

                    var skeletonGraphic = role.GetComponentInChildren<SkeletonGraphic>();
                    if (originMat == null)
                        originMat = skeletonGraphic.material;

                    skeletonGraphic.material = (isHave ? originMat : bookView.grayMat);

                    SpineUtil.Instance.PlayLoop(skeletonGraphic, isHave ? "stand" : null);
                }).Forget();

                // 名字
                JUIExtensions.txt(roleItemGo.transform.Find("NameText")).text = partnersCfg.Name;
                roleItemGo.SetActive(true);

                // 星级
                var startGo = roleItemGo.transform.Find("StarImg").gameObject;
                var startUpGo = roleItemGo.transform.Find("StartUpImg").gameObject;

                startGo.SetActive(isHave);
                startUpGo.SetActive(false);

                var star = 0;
                if (partnerInfo != null)
                    star = partnerInfo.Star;

                roleItemGo.transform.Find("StarImg/StartIconImg").gameObject.SetActive(star > 0);
                roleItemGo.transform.Find("StarImg/Text").GetComponent<Text>().text = (star > 0 ? star.ToString() : "");

                var upText = roleItemGo.transform.Find("UpText").GetComponent<Text>();
                upText.text = string.Empty;
                int nextStar = star == partnersCfg.StarConsume.Length ? star : star + 1;
                if (star < partnersCfg.StarAttribute.Length)
                {
                    var userdData = partnersCfg.GetStarConsumeByLv(nextStar);

                    canUp = true;
                    for (int j = 0; j < userdData.Length; j++)
                    {
                        if (!GameUtil.Instance.ComCoinToCheck(userdData[j].key, userdData[j].value, null, false))
                        {
                            canUp = false; break;
                        }
                    }

                    startUpGo.SetActive(canUp && isHave);

                    if (!isHave && canUp)
                        upText.text = "可合成";
                }
            }

            var role = Mathf.Ceil(partners.Count / 4f);
            var layoutElement = itemTrans.GetComponent<LayoutElement>();
            layoutElement.preferredHeight = role * 200 + (role - 1) * 10 + 42;
            itemTrans.gameObject.SetActive(true);
        }
    }

    private void OnRefresh_CommonResponse(GameEvent.Refrsh_Common _e)
    {
        if (isNullGo && _e.cid == CmdId.C2S_PartnerStar && _e.code != 0) return;

        if (roleTransDict.TryGetValue(selectItemId, out var trans))
        {
            var partnerInfo = PartnerM.Instance.GetOnePartnerInfoByCFGID(selectItemId);
            var canUp = false;
            var partnersCfg = CFG_Partner.GetById(partnerInfo.CfgId);

            var bgBtn = trans.Find("BgImg").GetComponent<Button>();
            bgBtn.onClick.RemoveAllListeners();
            bgBtn.onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                if (canUp && partnerInfo != null)
                    UITool.Open<UIPartnerUpStar>(null, partnerInfo.Uid, true);
                else
                    UITool.Open<UIStarDetails>(null, partnersCfg);
            });

            // 星级
            var startGo = trans.Find("StarImg").gameObject;
            var startUpGo = trans.Find("StartUpImg").gameObject;

            startGo.SetActive(true);
            startUpGo.SetActive(false);

            var star = 0;
            if (partnerInfo != null)
                star = partnerInfo.Star;

            trans.Find("StarImg/StartIconImg").gameObject.SetActive(star > 0);
            trans.Find("StarImg/Text").GetComponent<Text>().text = (star > 0 ? star.ToString() : "");

            int nextStar = star == partnersCfg.StarConsume.Length ? star : star + 1;
            if (star < partnersCfg.StarAttribute.Length)
            {
                var userdData = partnersCfg.GetStarConsumeByLv(nextStar);

                canUp = true;
                for (int j = 0; j < userdData.Length; j++)
                {
                    if (!GameUtil.Instance.ComCoinToCheck(userdData[j].key, userdData[j].value, null, false))
                    {
                        canUp = false; break;
                    }
                }

                startUpGo.SetActive(canUp);
            }
        }
    }

    public override void OnExit()
    {
        bookView = null;
        UniGameEvent.RemoveEventListener<GameEvent.Refrsh_Common>(OnRefresh_CommonResponse);
    }

    public override void OnDestroy()
    {
        bookView = null;
        UniGameEvent.RemoveEventListener<GameEvent.Refrsh_Common>(OnRefresh_CommonResponse);
    }
}