using HotUpdateScripts;
using HotUpdateScripts.Config;
using HotUpdateScripts.ProtoM;

using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;

using System;
using System.Collections.Generic;

using UnityEngine;
using UnityEngine.EventSystems;
using UnityEngine.UI;

using Object = UnityEngine.Object;

public class UIFightFormation : UIBase<UIModelFightFormation, UIViewFightFormation>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    public enum DragState
    {
        Begin,
        Dragging,
        End
    }

    public struct TeamInfo
    {
        public long uid;
        public string name;
        public bool isProtagonist;  // 是否是主角
        public int cfgId;
        public int level;
        public int quality;
        public int star;
        public int job;
        //public int power;
    }

    public override void OnInit()
    {
        GetView().btn_UnloadAll.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnUnloadAllTeams();
        });

        GetView().btn_UpLoadAll.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnAutoLoadAllTeams();
        });

        GetView().btn_Left.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            currentPage--;
            currentPage = Mathf.Max(currentPage, 1);
            RefreshAllTeams();
        });

        GetView().btn_Right.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            currentPage++;
            currentPage = Mathf.Min(currentPage, maxPage);
            RefreshAllTeams();
        });

        for (int i = 0; i < GetView().goTeams.Length; i++)
        {
            int index = i;
            var item = GetView().goTeams[i];

            item.btn_Remove.onClick.AddListener(() =>
            {
                AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
                OnRemoveTeam(index);
                RefreshAllTeams();
            });
            item.touchEvent.AddDragBeginListener((p) =>
            {
                OnDragGoTeam(index, DragState.Begin, p);
            });
            item.touchEvent.AddDragEndListener((p) =>
            {
                OnDragGoTeam(index, DragState.End, p);
            });
            item.touchEvent.AddDragListener((p) =>
            {
                OnDragGoTeam(index, DragState.Dragging, p);
            });
        }

        for (int i = 0; i < GetView().allTeams.Length; i++)
        {
            int index = i;
            var item = GetView().allTeams[i];

            item.touchEvent.AddDragBeginListener((p) =>
            {
                OnDragAlternativeTeam(index, DragState.Begin, p);
            });
            item.touchEvent.AddDragEndListener((p) =>
            {
                OnDragAlternativeTeam(index, DragState.End, p);
            });
            item.touchEvent.AddDragListener((p) =>
            {
                OnDragAlternativeTeam(index, DragState.Dragging, p);
            });
        }


        GetView().btn_Play.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            OnPlay();
        });

        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIFightFormation>();
        });
    }

    /// <summary>
    /// 要调战的人的uid
    /// </summary>
    long uid;

    long[] uids = new long[6];

    int pageCount = 12;
    int currentPage;
    int maxPage;
    bool isSendRequest;

    CFG_Room.RoomType roomType;

    private List<TeamInfo> allTeamsInfo = new List<TeamInfo>();

    private bool GetTeamInfo(long uid, out TeamInfo teamInfo)
    {
        int index = allTeamsInfo.FindIndex(team => team.uid == uid);
        if (index == -1)
            teamInfo = default;
        else
            teamInfo = allTeamsInfo[index];

        return index != -1;
    }

    private bool GetDisengagedTeam(out TeamInfo teamInfo)
    {
        foreach (var item in allTeamsInfo)
        {
            if (!Array.Exists(uids, uid => uid == item.uid))
            {
                teamInfo = item;
                return true;
            }
        }
        teamInfo = default;
        return false;
    }

    public override void OnEnter(params object[] args)
    {
        RegisterEvent();

        roomType = (CFG_Room.RoomType)args[0];
        isSendRequest = false;

        if (roomType == CFG_Room.RoomType.PVP_More)
        {
            uid = (long)args[1];
            Array.Copy(TurnBasedM.Instance.pvpFormationUids, uids, 6);
        }
        else if (roomType == CFG_Room.RoomType.WuXing)
        {
            Array.Copy(TurnBasedM.Instance.wuXingFormationUids, uids, 6);
        }

        currentPage = 1;

        LoadAllPartnets();
        RefreshGoTeams();
    }

    private void RegisterEvent()
    {
        UniGameEvent.AddEventListener<GameEvent.PVP_TurnBasedSlot>(Event_TurnBasedSlot);
        UniGameEvent.AddEventListener<GameEvent.WuXing_TurnBasedSlot>(Event_WuXingSlot);
        UniGameEvent.AddEventListener<GameEvent.TurnBasedChallenge>(Event_TurnBasedChallenge);
    }

    private void UnregisterEvent()
    {
        UniGameEvent.RemoveEventListener<GameEvent.PVP_TurnBasedSlot>(Event_TurnBasedSlot);
        UniGameEvent.RemoveEventListener<GameEvent.WuXing_TurnBasedSlot>(Event_WuXingSlot);
        UniGameEvent.RemoveEventListener<GameEvent.TurnBasedChallenge>(Event_TurnBasedChallenge);
    }

    /// <summary>
    /// 显示所有的伙伴，包括主角。
    /// </summary>
    private void LoadAllPartnets()
    {
        var partner = PartnerM.Instance.GetPartners();
        foreach (var item in partner)
        {
            var cfg_partner = CFG_Partner.GetById(item.CfgId);

            allTeamsInfo.Add(new TeamInfo
            {
                uid = item.Uid,
                name = CFG_Item.GetById(cfg_partner.ItemID).Name,
                isProtagonist = false,
                cfgId = item.CfgId,
                level = item.Level,
                quality = cfg_partner.Quality,
                star = item.Star,
                job = cfg_partner.job,
            });
        }

        allTeamsInfo.Sort((a, b) =>
        {
            if (a.quality != b.quality)
                return b.quality.CompareTo(a.quality);
            else if (a.level != b.level)
                return b.level.CompareTo(a.level);
            else if (a.star != b.star)
                return b.star.CompareTo(a.star);
            else
                return a.cfgId.CompareTo(b.cfgId);
        });

        allTeamsInfo.Insert(0, new TeamInfo
        {
            uid = RoleM.Instance.Uid,
            name = RoleM.Instance.Nickname,
            isProtagonist = true,
            cfgId = 0,
            level = RoleM.Instance.Level,
            quality = 0,
            star = 0,
            job = RoleM.Instance.Job,
        });

        maxPage = Mathf.CeilToInt((float)allTeamsInfo.Count / pageCount);

        RefreshAllTeams();
    }

    private void RefreshAllTeams()
    {
        int startIndex = (currentPage - 1) * pageCount;
        int endIndex = Mathf.Min(startIndex + pageCount - 1, allTeamsInfo.Count - 1);

        int currentPageCount = currentPage < maxPage ? 12 : allTeamsInfo.Count - (maxPage - 1) * pageCount;
        bool seatsIsFull = SeatsFull();

        ShowNavigationButton();

        for (int i = 0; i < 12; i++)
        {
            var item = GetView().allTeams[i];
            var dataIndex = startIndex + i;
            bool isShow = i < currentPageCount;
            item.Show(isShow);
            if (isShow)
            {
                var data = allTeamsInfo[dataIndex];

                if (data.isProtagonist)
                {
                    ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.RoleHeadAtlas, "img_" + RoleM.Instance.Job + RoleM.Instance.Sex, item.img_Head);
                    item.txt_Name.text = RoleM.Instance.Nickname;
                    item.outline_Name.effectColor = new Color(0.16f, 0.16f, 0.16f);
                    item.txt_Power.text = GameUtil.Instance.FormatPowerWithUnit(RoleM.Instance.Combat);
                    item.img_Power.enabled = true;
                }
                else
                {
                    var cfg_partnet = CFG_Partner.GetById(data.cfgId);
                    GameUtil.Instance.SetItemIcon(cfg_partnet.ItemID, item.img_Head);
                    item.txt_Name.text = GameUtil.Instance.GetItemName(cfg_partnet.ItemID);
                    item.outline_Name.effectColor = GameUtil.Instance.GetQualityOutlineColor(cfg_partnet.Quality);
                    var partnerStats = PartnerM.Instance.GetOnePartnerStats(data.uid);
                    item.img_Power.enabled = partnerStats != null;
                    if (partnerStats != null)
                        item.txt_Power.text = GameUtil.Instance.FormatPowerWithUnit(partnerStats.Combat);
                    else
                        item.txt_Power.text = "<color=#6f1c08>未出战</color>";
                }

                GameUtil.Instance.SetJobIcon(data.job % 10, item.img_Job);
                GameUtil.Instance.SetPartnetQualityFrame(data.quality, item.img_Quality);
                item.txt_Star.text = data.star.ToString();
                item.txt_Lv.text = data.level.ToString();
                item.txt_Star.gameObject.SetActive(!data.isProtagonist);
                int findIndex = Array.FindIndex(uids, uid => uid == data.uid);
                if (findIndex != -1)
                {
                    item.obj_Mask.SetActive(true);
                    item.txt_MaskTip.text = $"{findIndex + 1}阵位";
                }
                else if (seatsIsFull)
                {
                    item.obj_Mask.SetActive(true);
                    item.txt_MaskTip.text = "阵位\n已满";
                }
                else
                {
                    item.obj_Mask.SetActive(false);
                }
            }
        }
    }

    /// <summary>
    /// 显示隐藏左右换页按钮
    /// </summary>
    private void ShowNavigationButton()
    {
        GetView().btn_Left.gameObject.SetActive(currentPage > 1);
        GetView().btn_Right.gameObject.SetActive(currentPage < maxPage);
    }

    private void RefreshGoTeams()
    {
        for (int i = 0; i < uids.Length; i++)
        {
            if (uids[i] != 0)
            {
                if (GetTeamInfo(uids[i], out TeamInfo teamInfo))
                {
                    SetItemAndShow(GetView().goTeams[i], true, teamInfo);
                }
                else
                {
                    SetItemAndShow(GetView().goTeams[i], false);
                }
            }
            else
            {
                SetItemAndShow(GetView().goTeams[i], false);
            }
        }
    }

    private void SetItemAndShow(UIViewFightFormation.UpTeam upTeam, bool isShow, TeamInfo teamInfo = default)
    {
        upTeam.ShowHead(isShow);
        upTeam.ShowPositionText(!isShow);

        if (isShow)
        {
            bool isProtagonist = teamInfo.uid == RoleM.Instance.Uid;
            if (isProtagonist)
            {
                ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.RoleHeadAtlas, "img_" + RoleM.Instance.Job + RoleM.Instance.Sex, upTeam.img_Head);
            }
            else
            {
                var cfg_partnet = CFG_Partner.GetById(teamInfo.cfgId);
                GameUtil.Instance.SetItemIcon(cfg_partnet.ItemID, upTeam.img_Head);
                upTeam.txt_Star.text = teamInfo.star.ToString();
            }

            GameUtil.Instance.SetJobIcon(teamInfo.job % 10, upTeam.img_Job);
            GameUtil.Instance.SetPartnetQualityFrame(teamInfo.quality, upTeam.img_Quality);
            upTeam.txt_Lv.text = teamInfo.level.ToString();
            upTeam.txt_Star.gameObject.SetActive(!isProtagonist);
        }
    }

    private void OnUnloadAllTeams()
    {
        for (int i = 0; i < uids.Length; i++)
        {
            if (uids[i] != 0)
                OnRemoveTeam(i);
        }

        RefreshAllTeams();
    }

    private void OnAutoLoadAllTeams()
    {
        for (int i = 0; i < uids.Length; i++)
        {
            if(SeatsFull())
            {
                break;
            }
            else if (uids[i] == 0)
            {
                if (GetDisengagedTeam(out TeamInfo teamInfo))
                    OnGoTeam(i, teamInfo.uid);
                else // 如果没有可上阵的角色，直接退出
                    break;
            }
        }

        RefreshAllTeams();
    }

    // 拖拽的已上阵的角色的Item拷贝
    private RectTransform dragMoveGoItem;

    // 已上阵角色的拖拽
    private void OnDragGoTeam(int position, DragState dragState, PointerEventData pointerEventData)
    {
        int _position = position;
        if (uids[position] == 0)
            return;

        if (dragState == DragState.Begin)
        {
            dragMoveGoItem = Object.Instantiate(GetView().goTeams[position].head, GetView().root_Panel).GetComponent<RectTransform>();
            dragMoveGoItem.anchorMax = new Vector2(0.5f, 0.5f);
            dragMoveGoItem.anchorMin = new Vector2(0.5f, 0.5f);
            dragMoveGoItem.transform.Find("btn_remove").gameObject.SetActive(false);
            GetView().goTeams[position].ShowHead(false);
            GetView().goTeams[position].ShowPositionText(true);

            var imgs = dragMoveGoItem.GetComponentsInChildren<Image>();
            foreach (Image image in imgs)
            {
                image.raycastTarget = false;
            }
        }
        else if (dragState == DragState.Dragging)
        {
            RectTransformUtility.ScreenPointToLocalPointInRectangle(GetView().root_Panel.GetComponent<RectTransform>(), Input.mousePosition, JUIExtensions.uiCamera, out Vector2 localPoint);
            dragMoveGoItem.anchoredPosition = localPoint;
        }
        else if (dragState == DragState.End)
        {
            if (Object.Equals(pointerEventData.pointerCurrentRaycast.gameObject, null))
            {
                GetView().goTeams[position].ShowHead(true);
                GetView().goTeams[position].ShowPositionText(false);
                Object.Destroy(dragMoveGoItem.gameObject);
                dragMoveGoItem = null;
                return;
            }

            string objName = pointerEventData.pointerCurrentRaycast.gameObject.name;
            if (objName == "btn_remove") // 此位置已被占用
            {
                var target = pointerEventData.pointerCurrentRaycast.gameObject.transform.parent.parent;
                int _targetPosition = target.GetSiblingIndex();
                // 正在拖拽的uid
                long opUid = uids[_position];
                // 将正在拖拽的角色与目标角色互换位置
                OnGoTeam(_position, uids[_targetPosition]);
                OnGoTeam(_targetPosition, opUid);
            }
            else if (pointerEventData.pointerCurrentRaycast.gameObject.name.Contains("PosItem_"))
            {
                var target = pointerEventData.pointerCurrentRaycast.gameObject.transform;
                int _targetPosition = target.GetSiblingIndex();
                long opUid = uids[_position];

                if (uids[_targetPosition] != 0)
                {
                    // 将正在拖拽的角色与目标角色互换位置
                    OnGoTeam(_position, uids[_targetPosition]);
                    OnGoTeam(_targetPosition, opUid);
                    RefreshAllTeams();
                }
                else
                {
                    // 目标位置没有角色，直接上阵，并腾出原位置
                    OnGoTeam(_targetPosition, opUid);
                    OnRemoveTeam(_position);
                    RefreshAllTeams();
                }
            }
            else
            {
                GetView().goTeams[position].ShowHead(true);
                GetView().goTeams[position].ShowPositionText(false);
            }

            Object.Destroy(dragMoveGoItem.gameObject);
            dragMoveGoItem = null;
        }
    }

    // 拖拽的可供上阵的角色列表的Item拷贝
    private RectTransform dragMoveAlternativeItem;

    // 可选角色列表的拖拽
    private void OnDragAlternativeTeam(int position, DragState dragState, PointerEventData pointerEventData)
    {
        int _position = position;
        int dataIndex = GetDataIndex(_position);

        if (SeatsFull() || Array.Exists(uids, uid => uid == allTeamsInfo[dataIndex].uid))
        {
            return;
        }

        if (dragState == DragState.Begin)
        {
            dragMoveAlternativeItem = Object.Instantiate(GetView().allTeams[_position].go, GetView().root_Panel).GetComponent<RectTransform>();
            dragMoveAlternativeItem.anchorMax = new Vector2(0.5f, 0.5f);
            dragMoveAlternativeItem.anchorMin = new Vector2(0.5f, 0.5f);
            dragMoveAlternativeItem.Find("img_fight").gameObject.SetActive(false);

            var imgs = dragMoveAlternativeItem.GetComponentsInChildren<Image>();
            foreach (Image image in imgs)
            {
                image.raycastTarget = false;
            }
        }
        else if (dragState == DragState.Dragging)
        {
            RectTransformUtility.ScreenPointToLocalPointInRectangle(GetView().root_Panel.GetComponent<RectTransform>(), Input.mousePosition, JUIExtensions.uiCamera, out Vector2 localPoint);
            dragMoveAlternativeItem.anchoredPosition = localPoint;
        }
        else if (dragState == DragState.End)
        {
            if (Object.Equals(pointerEventData.pointerCurrentRaycast.gameObject, null))
            {
                Object.Destroy(dragMoveAlternativeItem.gameObject);
                dragMoveAlternativeItem = null;
                return;
            }

            string objName = pointerEventData.pointerCurrentRaycast.gameObject.name;
            if (objName == "btn_remove")
            {
                var target = pointerEventData.pointerCurrentRaycast.gameObject.transform.parent.parent;
                int _targetPosition = target.GetSiblingIndex();
                // 替换已上阵的角色
                OnGoTeam(_targetPosition, allTeamsInfo[dataIndex].uid);
                RefreshAllTeams();
            }
            else if (pointerEventData.pointerCurrentRaycast.gameObject.name.Contains("PosItem_"))
            {
                var target = pointerEventData.pointerCurrentRaycast.gameObject.transform;
                int _targetPosition = target.GetSiblingIndex();
                if (uids[_targetPosition] == 0 && SeatsFull())
                {
                    Object.Destroy(dragMoveAlternativeItem.gameObject);
                    dragMoveAlternativeItem = null;
                    UITool.ShowTips("阵容已满！");
                    return;
                }
                OnGoTeam(_targetPosition, allTeamsInfo[dataIndex].uid);
                RefreshAllTeams();
            }
            Object.Destroy(dragMoveAlternativeItem.gameObject);
            dragMoveAlternativeItem = null;
        }
    }

    /// <summary>
    /// 下阵
    /// </summary>
    /// <param name="pos">要下阵的位置</param>
    private void OnRemoveTeam(int pos)
    {
        uids[pos] = 0;
        SetItemAndShow(GetView().goTeams[pos], false);
    }

    /// <summary>
    /// 上阵或替换
    /// </summary>
    /// <param name="pos">要上阵的位置</param>
    /// <param name="uid">要上阵的角色uid</param>
    private void OnGoTeam(int pos, long uid)
    {
        uids[pos] = uid;
        GetTeamInfo(uid, out TeamInfo teamInfo);
        SetItemAndShow(GetView().goTeams[pos], true, teamInfo);
    }

    private void OnPlay()
    {
        if (isSendRequest) return;

        if (!CheckFormation()) return;

        if(roomType == CFG_Room.RoomType.PVP_More)
        {
            TurnBasedM.Instance.opponentUid = uid;
            TcpC.C2S_TurnBasedSlotRequest(uids);
            Debug.Log("---- <color=cyan>回合制 pvp: 保存队形！</color>");
            isSendRequest = true;
            UITool.Open<UILoading>(null, new UILoading.Loading { _time = 0 });

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "玩法",
                component = "回合挑战",
                position = "点击进入选择对手",
                event_click = $"布阵开始战斗",
            });
        }
        else if(roomType == CFG_Room.RoomType.WuXing)
        {
            TurnBasedM.Instance.currentChallengeLayer = TurnBasedM.Instance.elementScore == 0 ?
                1 : CFG_ElementTrials.GetLayerById(TurnBasedM.Instance.elementScore) + 1;

            TcpC.C2S_ElementPartnerSlotRequest(uids);
            Debug.Log("---- <color=yellow>回合制 五行: 保存队形！</color>");
            isSendRequest = true;
            UITool.Open<UILoading>(null, new UILoading.Loading { _time = 0 });

            HttpC.ShuShuSendClickEvent(new HTTP_WindowShushuClickEventRequest
            {
                page = "五行试炼界面",
                component = "点击挑战",
                position = "布阵界面",
                event_click = $"开始战斗",
            });
        }
        else
        {
            return;
        }

    }

    private bool CheckFormation()
    {
        // 检查是否有上阵的角色
        bool isEmpty = true;
        // 检查是否有主角
        bool haveProtagonist = false;

        for (int i = 0; i < uids.Length; i++)
        {
            if (isEmpty && uids[i] != 0)
            {
                isEmpty = false;
            }

            if (!haveProtagonist && uids[i] == RoleM.Instance.Uid)
            {
                haveProtagonist = true;
            }

            if (!isEmpty && haveProtagonist) break;
        }

        if (isEmpty)
        {
            UITool.Open<UIMessageBox>(null, new UIMessageBox.CallBack
            {
                ShowStr = "当前上阵为空，是否需要系统自动分配阵位？",
                BtnNo = "手动上阵",
                BtnYes = "自动上阵",
                YesFun = OnAutoLoadAllTeams,
            });
            return false;
        }

        if (!haveProtagonist)
        {
            UITool.ShowTips("请上阵主角后再进行挑战！");
            return false;
        }

        return true;
    }

    private int GetDataIndex(int position)
    {
        int startIndex = (currentPage - 1) * pageCount;
        return startIndex + position;
    }

    /// <summary>
    /// 阵容位置是否已满
    /// </summary>
    /// <returns></returns>
    private bool SeatsFull()
    {
        int emptyCount = 0;
        for (int i = 0; i < uids.Length; i++)
        {
            if (uids[i] == 0)
            {
                emptyCount++;
            }
        }

        return emptyCount <= 1;
    }

    /// <summary>
    /// 阵容设置。
    /// </summary>`
    /// <param name="_e"></param>
    private void Event_TurnBasedSlot(GameEvent.PVP_TurnBasedSlot _e)
    {
        // 阵容设置成功
        Debug.Log("----- <color=yellow>pvp： 队形保存成功。。。</color>");
        Array.Copy(TurnBasedM.Instance.pvpFormationUids, uids, 6);
        TcpC.C2S_TurnBasedChallengeRequest(uid);
        isSendRequest = false;
    }

    private void Event_WuXingSlot(GameEvent.WuXing_TurnBasedSlot _e)
    {
        Debug.Log("----- <color=yellow>五行： 队形保存成功。。。</color>");
        // 五行试炼阵容设置成功
        Array.Copy(TurnBasedM.Instance.wuXingFormationUids, uids, 6);
        TcpC.C2S_ElementChallengeRequest();
        isSendRequest = false;
    }

    /// <summary>
    /// 发起挑战
    /// </summary>
    /// <param name="_e"></param>
    private void Event_TurnBasedChallenge(GameEvent.TurnBasedChallenge _e)
    {
        // 开始挑战
        UITool.Close<UIFightFormation>();
    }

    public override void OnExit()
    {
        UnregisterEvent();

        for (int i = 0; i < uids.Length; i++)
        {
            uids[i] = 0;
        }
        allTeamsInfo.Clear();
    }

    public override void OnDestroy()
    {

    }
}