using Cysharp.Threading.Tasks;

using DG.Tweening;

using HotUpdateScripts;
using HotUpdateScripts.Game;
using HotUpdateScripts.ProtoM;

using JEngine.Core;
using JEngine.Helper;
using JEngine.UI;

using MetaFramework.UI;

using pb;

using System;
using System.Collections.Generic;
using System.Text;

using UnityEngine;
using UnityEngine.UI;


public class UIGame3D : UIBase<UIModelGame3D, UIViewGame3D>
{
    public override UILayer GetLayer()
    {
        return UILayer.Base;
    }
    public sealed class HUD
    {
        public GameObject obj;
        public Slider HBslider;
        public Text text;
        public Text text_HP;
        public GameObject obj_Name;
        public GameObject baseObj;
        public Image Img_Fill;
        public Image Img_NextFile;
        public Image[] Img_buffers;
        public bool IsUsed;
        public Button _btn;
        public Text _txt_GuiSu;
        public Image img_title;
    }
    private static Dictionary<long, HUD> HuDDic = new Dictionary<long, HUD>();

    private string[] BossHubColor = { "#2a76c4", "#5cc42a", "#d6c730", "#c4692a", "#c42a2a" };

    private Dictionary<HUB, GameObject> HubToGoDic = new Dictionary<HUB, GameObject>();

    private Dictionary<long, GameObject> TipEnemyDic = new Dictionary<long, GameObject>();
    private Queue<GameObject> TipEnemyQueue = new Queue<GameObject>();


    private bool isShowEmenyTip = true;

    private UIViewGame3D _UIViewGame3D;
    public override void OnInit()
    {
        _UIViewGame3D = GetView();

        HubToGoDic.Add(HUB.BossHUB, _UIViewGame3D.root_BossMonterHub.gameObject);
        HubToGoDic.Add(HUB.PlayerHUB, _UIViewGame3D.obj_PlayerHub);
        HubToGoDic.Add(HUB.MonsterHUB, _UIViewGame3D.obj_Monster);
        HubToGoDic.Add(HUB.CollectHUB, _UIViewGame3D.root_Collect);

        _UIViewGame3D.UpateCallBack = SetUUDShow;

        isShowEmenyTip = GameRoomM.Instance.CFGRoom.Roomtype != CFG_Room.RoomType.PVE_NoPeace;

        GetView().root_Pos.SetActive(GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVE_NoPeace);
        if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVE_NoPeace)
            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.FightAtlas, GameRoomM.Instance.CFGRoom.Name.Contains("地狱") ? "bg_qxzc_map_2" : "bg_qxzc_map_1", GetView().root_MapBG.GetComponent<Image>());

        GetView().root_Map.GetComponent<Button>().onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Open<UIGameFightWarMap>();
        });
    }

    public enum HUB
    {
        PlayerHUB,
        MonsterHUB,
        BossHUB,
        CollectHUB,
    }

    private void GetEnemyTip(long id, Vector2 uiPos)
    {
        if (!isShowEmenyTip)
            return;

        if (!TipEnemyDic.ContainsKey(id))
        {
            GameObject go;
            if (TipEnemyQueue.Count <= 0)
            {
                go = GameObject.Instantiate(_UIViewGame3D.tip_enemy, _UIViewGame3D.rect_Monster.transform);
                go.transform.SetParent(_UIViewGame3D.rect_Monster);
                go.transform.localScale = Vector3.one;
            }
            else
            {
                go = TipEnemyQueue.Dequeue();
            }
            go.SetActive(true);
            TipEnemyDic.Add(id, go);

        }

        if (Mathf.Abs(uiPos.x) > -_UIViewGame3D.rect_Monster.rect.x)
        {
            uiPos.x = uiPos.x > 0 ? _UIViewGame3D.rect_Monster.rect.xMax - 20 : _UIViewGame3D.rect_Monster.rect.xMin + 20;
        }

        if (Mathf.Abs(uiPos.y) > -_UIViewGame3D.rect_Monster.rect.y)
        {
            uiPos.y = uiPos.y > 0 ? _UIViewGame3D.rect_Monster.rect.yMax : _UIViewGame3D.rect_Monster.rect.yMin + 20;
        }
        float offy = (uiPos.y - _UIViewGame3D.rect_Monster.rect.y) / _UIViewGame3D.rect_Monster.rect.size.y;
        float offx = (uiPos.x - _UIViewGame3D.rect_Monster.rect.x) / _UIViewGame3D.rect_Monster.rect.size.x;

        Vector3 rot;
        if (offy < 0.3333)
        {
            if (offx < 0.3333)
                rot = new Vector3(0, 0, 135);
            else if (offx < 0.6666)
                rot = new Vector3(0, 0, 180);
            else
                rot = new Vector3(0, 0, 225);
        }
        else if (offy < 0.6666)
        {
            if (offx < 0.3333)
                rot = new Vector3(0, 0, 90);
            else
                rot = new Vector3(0, 0, 270);
        }
        else
        {
            if (offx < 0.3333)
                rot = new Vector3(0, 0, 45);
            else if (offx < 0.6666)
                rot = new Vector3(0, 0, 0);
            else
                rot = new Vector3(0, 0, 315);
        }

        if (ShouldUpdate(TipEnemyDic[id].transform.localPosition, uiPos))
        {
            TipEnemyDic[id].transform.localPosition = uiPos;
            TipEnemyDic[id].transform.localEulerAngles = rot;
        }
    }


    private void RecycEnemyTip(long uid)
    {
        if (TipEnemyDic.ContainsKey(uid))
        {
            TipEnemyDic[uid].SetActive(false);
            TipEnemyQueue.Enqueue(TipEnemyDic[uid]);
            TipEnemyDic.Remove(uid);
        }
    }

    private HUD GetHUDById(long id, HUB hub = HUB.PlayerHUB)
    {
        if (!HuDDic.ContainsKey(id))
        {
            string Key = hub.ToString();
            GameObject obj = HubToGoDic[hub];
            if (!PoolHelper.Instance.CheckPool(Key))
            {
                PoolHelper.Instance.AddPool(Key, obj);
            }
            GameObject go = PoolHelper.Instance.GetPoolObj(Key);
            go.name = Key;
            HuDDic.Add(id, new HUD
            {
                obj = go,
                HBslider = go.GetComponent<Slider>(),
                Img_Fill = JUIExtensions.img(go.transform.Find("Fill Area/Fill")),
                IsUsed = false,
            });

            if (go.transform.Find("txt_name"))
            {
                HuDDic[id].text = JUIExtensions.txt(go.transform.Find("txt_name"));
            }

            if (go.transform.Find("img_title"))
            {
                HuDDic[id].img_title = JUIExtensions.img(go.transform.Find("img_title"));
            }

            if (go.transform.Find("root"))
            {
                HuDDic[id]._txt_GuiSu = JUIExtensions.txt(go.transform.Find("root/imgBg/txt_name"));
                HuDDic[id].obj_Name = go.transform.Find("root/imgBg").gameObject;
                HuDDic[id].baseObj = go.transform.Find("root").gameObject;
                HuDDic[id]._btn = go.transform.Find("root/btn").GetComponent<Button>();
            }

            if (go.transform.Find("txt_HP"))
            {
                HuDDic[id].text_HP = JUIExtensions.txt(go.transform.Find("txt_HP"));
            }


            if (go.transform.Find("NextFile"))
            {
                HuDDic[id].Img_NextFile = JUIExtensions.img(go.transform.Find("NextFile"));
            }

            if (go.transform.Find("NextFile"))
            {
                HuDDic[id].Img_NextFile = JUIExtensions.img(go.transform.Find("NextFile"));
            }

            if (go.transform.Find("bufferRoot"))
            {
                var bufferRoot = go.transform.Find("bufferRoot");
                HuDDic[id].Img_buffers = new Image[bufferRoot.childCount];
                for (int i = 0; i < bufferRoot.childCount; i++)
                {
                    HuDDic[id].Img_buffers[i] = JUIExtensions.img(bufferRoot.GetChild(i));
                }
            }


            go.transform.SetParent(hub == HUB.PlayerHUB ? _UIViewGame3D.root_PlayerHub : _UIViewGame3D.root_Monter);
            go.transform.localScale = Vector3.one;
            go.SetActive(true);
        }
        return HuDDic[id];
    }

    private void SetHUDById(long id)
    {
        if (HuDDic.ContainsKey(id))
        {
            PoolHelper.Instance.RecycObj(HuDDic[id].obj);
            HuDDic.Remove(id);
        }
    }

    private void CleatHUB()
    {
        if (PoolHelper.Instance == null)
            return;
        foreach (var item in HuDDic)
        {
            PoolHelper.Instance.RecycObj(item.Value.obj);
        }
        HuDDic.Clear();
    }

    public override void OnEnter(params object[] args)
    {
        OnEventAddOrReMove(true);

        if (GuidanceM.Instance.CheckShowGuide() || GuidanceM.Instance.isGuiding)
        {
            GuidanceM.Instance.isGuiding = true;
            UITool.Open<UIGuide>();
        }
    }

    bool ShouldUpdate(Vector2 newPos, Vector2 oldPos)
    {
        return Vector2.Distance(newPos, oldPos) > 1f;
    }

    public void SetUUDShow()
    {
        if (isNullGo || GameRoomM.Instance._GameState == GameRoomM.GameState.Over || GameRoomM.Instance._GameState == GameRoomM.GameState.No)
            return;

        // 缓存常用组件
        var monsterRect = _UIViewGame3D.rect_Monster;
        var rect = monsterRect.rect;
        float rectXMax = rect.xMax;
        float rectXMin = rect.xMin;
        float rectYMax = rect.yMax;
        float rectYMin = rect.yMin;

        // 优化1: 缓存怪物列表
        var monsters = MonsterManager.Instance.getMonsterObjs();
        int monsterCount = monsters.Count;

        // 怪物处理
        for (int i = 0; i < monsterCount; i++)
        {
            var monsterObj = monsters[i];
            if (monsterObj?._monser == null) continue;

            var monster = monsterObj._monser;
            long monsterId = monster.uID;
            Vector3 hubPos = monster.root_hub.transform.position;
            Vector2 uiPos = GetCanvarPos(hubPos);

            bool isOutside = Mathf.Abs(uiPos.x) > -rectXMin || Mathf.Abs(uiPos.y) > -rectYMin;

            if (isOutside)
            {
                SetHUDById(monsterId);

                GetEnemyTip(monsterId, uiPos);
            }
            else
            {
                RecycEnemyTip(monsterId);

                var cfgMonster = CFG_Monster.GetById(monster.cfgId);
                bool isBoss = cfgMonster.Isboss == 1 || cfgMonster.Isboss == 2 || cfgMonster.Isboss == 4;
                bool isCollect = cfgMonster.Isboss == 3;
                HUB hubType = isBoss ? HUB.BossHUB :
                              isCollect ? HUB.CollectHUB : HUB.MonsterHUB;

                var hud = GetHUDById(monsterId, hubType);
                Transform hudTransform = hud.obj.transform;

                if (ShouldUpdate(hudTransform.localPosition, uiPos))
                    hudTransform.localPosition = uiPos;

                if (isCollect)
                {
                    hud.baseObj.SetActive(monster.curLifeTime == 0);
                }

                if (!hud.IsUsed)
                {
                    hud.IsUsed = true;

                    if (isCollect)
                    {
                        bool isShowBtn = monster._Ownership == null || monster._Ownership.Type == OwnershipType.OWNERSHIP_NONE || (monster._Ownership.Type == OwnershipType.OWNERSHIP_PLAYER && monster._Ownership.Uid == MainRole._MainRole.uID) || (monster._Ownership.Type == OwnershipType.OWNERSHIP_GUILD && monster._Ownership.Uid == MainRole._MainRole.GuildId);
                        hud.obj_Name.SetActive(!isShowBtn);
                        if (!isShowBtn)
                        {
                            hud._txt_GuiSu.text = monster._Ownership.Type == OwnershipType.OWNERSHIP_NONE ? "无归属" : $"{(monster._Ownership.Type == OwnershipType.OWNERSHIP_PLAYER ? "玩家" : "行会")}({monster._Ownership.Name})归属";
                        }
                        hud._btn.gameObject.SetActive(isShowBtn);
                    }

                    hud.text.text = CFG_Monster.GetHUDNameById(monster.cfgId);

                    if (isCollect && hud._btn != null)
                    {
                        hud._btn.onClick.RemoveAllListeners();
                        hud._btn.onClick.AddListener(() =>
                        {
                            if (GameRoomM.Instance.CollectMonsterDic.ContainsKey(monsterId))
                            {
                                UITool.ShowTips($"玩家{GameRoomM.Instance.CollectMonsterDic[monsterId]}正在开启");
                                return;
                            }

                            if (monster._Ownership != null)
                            {
                                if ((monster._Ownership.Type == OwnershipType.OWNERSHIP_PLAYER && monster._Ownership.Uid != MainRole._MainRole.uID)
                                    || (monster._Ownership.Type == OwnershipType.OWNERSHIP_GUILD && monster._Ownership.Uid != MainRole._MainRole.GuildId))
                                {
                                    UITool.ShowTips($"该物品不属于你！");
                                    return;
                                }
                            }

                            MainRole._MainRole.targetCollectID = monsterId;
                        });
                    }
                }

                // 优化2: 使用局部变量代替多次访问
                int maxHp = monster.maxHp;
                int curHp = monster.hp;

                switch (hubType)
                {
                    case HUB.MonsterHUB:
                        hud.HBslider.value = (float)curHp / maxHp;
                        break;

                    case HUB.CollectHUB:
                        break;

                    case HUB.BossHUB:
                        int oneHP = maxHp / (int)cfgMonster.Barnum;
                        int curHpMod = curHp == 0 ? 0 : (curHp % oneHP == 0 ? oneHP : curHp % oneHP);
                        int curCount = Mathf.CeilToInt((float)curHp / oneHP);

                        hud.HBslider.value = (float)curHpMod / oneHP;
                        hud.text_HP.text = "X" + curCount;

                        // 优化3: 避免数组越界
                        int colorIndex = curCount % BossHubColor.Length;
                        hud.Img_Fill.color = JUIExtensions.ToColor(BossHubColor[colorIndex]);

                        bool showNext = curCount > 1;
                        hud.Img_NextFile.enabled = showNext;
                        if (showNext)
                        {
                            int prevColorIndex = (curCount - 1) % BossHubColor.Length;
                            hud.Img_NextFile.color = JUIExtensions.ToColor(BossHubColor[prevColorIndex]);
                        }
                        break;
                }

                // 优化4: 提前检查避免不必要操作
                if (hubType != HUB.CollectHUB && hud.Img_buffers != null)
                {
                    var bufferList = monster._BufferList;
                    int bufferCount = Mathf.Min(hud.Img_buffers.Length, bufferList.Count);

                    for (int j = 0; j < hud.Img_buffers.Length; j++)
                    {
                        bool shouldShow = j < bufferCount;
                        var imgBuffer = hud.Img_buffers[j];
                        imgBuffer.gameObject.SetActive(shouldShow);

                        if (shouldShow)
                        {
                            int buffId = bufferList[j];
                            var cfgBuff = CFG_BuffEffect.GetById(buffId);
                            if (!string.IsNullOrEmpty(cfgBuff.Icon))
                            {
                                ResouseHelper.Instance.LoadSpriteByAtlas(
                                    cfgBuff.Gallery,
                                    cfgBuff.Icon,
                                    imgBuffer
                                );
                            }
                        }
                    }
                }
            }
        }

        // 玩家处理 - 优化逻辑结构
        if (GameRoomM.Instance.IsTurnBasedRoom())
            ProcessTurnBasedPlayers(_UIViewGame3D, monsterRect, rectXMin, rectYMin);
        else
        {
            ProcessRealTimePlayers(_UIViewGame3D, monsterRect, rectXMin, rectYMin);
            OnCllectSliderShow();

            OnRefreshMap();
        }

    }

    private void OnCllectSliderShow()
    {
        bool isShow = false;
        if (MainRole._MainRole != null)
        {
            var monset = MonsterManager.Instance.getOneMonster(MainRole._MainRole.targetCollectID);
            if (monset != null && monset._monser.curLifeTime > 0)
            {
                isShow = true;
                _UIViewGame3D.slider_Collect.value = monset._monser.curLifeTime / monset._monser.maxLifeTime;
                _UIViewGame3D.txt_SiliderCollect.text = $"开启中{(monset._monser.maxLifeTime - monset._monser.curLifeTime).ToString("0.0")}S";
            }
        }
        _UIViewGame3D.slider_Collect.gameObject.SetActive(isShow);
    }

    // 提取出来的方法 - 回合制玩家
    private void ProcessTurnBasedPlayers(UIViewGame3D view, RectTransform rect, float rectXMin, float rectYMin)
    {
        var players = TurnBasedRoleManager.Instance.allRoles;
        int playerCount = players.Count;
        bool isPlayerAttack = GameRoomM.Instance.GameAttack == GameRoomM.AttackClass.Player;

        for (int i = 0; i < playerCount; i++)
        {
            var role = players[i];
            if (role == null) continue;

            Vector2 uiPos = GetCanvarPos(role.root_hub.position);
            bool isOutside = Mathf.Abs(uiPos.x) > -rectXMin || Mathf.Abs(uiPos.y) > -rectYMin;
            bool isDead = role.currentHp <= 0;

            if (isDead || isOutside)
            {
                if (isOutside) SetHUDById(role.uid);

                if (isPlayerAttack && isOutside && !isDead)
                    GetEnemyTip(role.uid, uiPos);
                else
                    RecycEnemyTip(role.uid);
            }
            else
            {
                RecycEnemyTip(role.uid);
                var hud = GetHUDById(role.uid, HUB.PlayerHUB);

                hud.obj.transform.localPosition = uiPos;
                hud.HBslider.value = (float)role.currentHp / role.maxHp;

                string name = role.roleType switch
                {
                    ObjectType.PLAYER => GetRoleStr(GetPlayerNameWithGuild(role.uid), role.uid == RoleM.Instance.Uid, 0),
                    ObjectType.PARTNER => GameUtil.Instance.GetStrByQuality(CFG_Partner.GetById(role.cfgId).Quality, role.roleName),
                    ObjectType.MONSTER => CFG_Monster.GetHUDNameById(role.cfgId),
                    _ => ""
                };

                if (!hud.IsUsed)
                {
                    hud.IsUsed = true;
                    string nameValue = null;

                    if (TurnBasedM.Instance.isOpen)
                    {
                        if (role.campType == TurnBasedBattle.ActionCamp.BLUE_CAMP)
                        {
                            if (role.roleType == ObjectType.PLAYER)
                                nameValue = $"{name}{GetRoleStr(TurnBasedM.Instance.GetElementName(role.mainElement, true), role.uid == RoleM.Instance.Uid, 0)}";
                            else if (role.roleType == ObjectType.PARTNER)
                                nameValue = $"{name}{GameUtil.Instance.GetStrByQuality(CFG_Partner.GetById(role.cfgId).Quality, TurnBasedM.Instance.GetElementName(role.mainElement, true))}";
                        }
                        else if (role.campType == TurnBasedBattle.ActionCamp.RED_CAMP)
                        {
                            if (role.roleType == ObjectType.PLAYER)
                                nameValue = $"{GetRoleStr(TurnBasedM.Instance.GetElementName(role.mainElement, true), role.uid == RoleM.Instance.Uid, 0)}{name}";
                            else if (role.roleType == ObjectType.PARTNER)
                                nameValue = $"{GameUtil.Instance.GetStrByQuality(CFG_Partner.GetById(role.cfgId).Quality, TurnBasedM.Instance.GetElementName(role.mainElement, true))}{name}";
                            else
                                nameValue = name;
                        }
                    }
                    else
                    {
                        nameValue = name;
                    }

                    if (hud.img_title != null)
                        GameUtil.Instance.SetTitleIcon(0, hud.img_title);
                    hud.text.text = nameValue;
                }

                // 优化buff处理
                ProcessBuffers(hud, role.bufferList);
            }
        }

        static string GetPlayerNameWithGuild(long uid)
        {
            var roleInfo = TurnBasedM.Instance.GetRoleInfo(uid);
            if (roleInfo == null) return null;

            if (string.IsNullOrEmpty(roleInfo.GuildName))
                return roleInfo.Name;
            else
                return $"[{roleInfo.GuildName}]\n{roleInfo.Name}";
        }
    }


    // 提取出来的方法 - 实时玩家
    private void ProcessRealTimePlayers(UIViewGame3D view, RectTransform rect, float rectXMin, float rectYMin)
    {
        var players = FightController.Instance.GetRolesInfo();
        int playerCount = players.Count;
        bool isPlayerAttack = GameRoomM.Instance.GameAttack == GameRoomM.AttackClass.Player;

        for (int i = 0; i < playerCount; i++)
        {
            var playerInfo = players[i];
            if (playerInfo?.role == null) continue;

            var role = playerInfo.role;
            if (role._ObjectType != pb.ObjectType.PLAYER &&
                role._ObjectType != pb.ObjectType.PARTNER) continue;

            Vector2 uiPos = GetCanvarPos(role.root_hub.position);
            bool isOutside = Mathf.Abs(uiPos.x) > -rectXMin || Mathf.Abs(uiPos.y) > -rectYMin;
            bool isDead = role.hp <= 0;

            long uid = role.uID;
            //安全区显示
            if (uid == MainRole._MainRole.uID && _UIViewGame3D.root_Pos.activeSelf)
            {
                bool isSafe = SceneManager.IsSafePos(role._gridNow.X, role._gridNow.Y);
                _UIViewGame3D.txt_Pos.text = $"[{role._gridNow.X},{role._gridNow.Y}]{(isSafe ? "<color=#4dae1f>安全区</color>" : "<color=#ff6868>战斗区</color>")}";
            }

            if (isDead || isOutside)
            {
                SetHUDById(uid);

                if (isPlayerAttack && isOutside && !isDead)
                    GetEnemyTip(uid, uiPos);
                else
                    RecycEnemyTip(uid);
            }
            else
            {
                RecycEnemyTip(uid);
                var hud = GetHUDById(uid, HUB.PlayerHUB);

                if (ShouldUpdate(hud.obj.transform.localPosition, uiPos))
                    hud.obj.transform.localPosition = uiPos;

                hud.HBslider.value = (float)role.hp / role.maxHp;

                if (!hud.IsUsed)
                {
                    hud.IsUsed = true;
                    // 优化字符串生成
                    string displayName;
                    if (role._ObjectType == pb.ObjectType.PARTNER && role.OwnerUid != 0)
                    {
                        string ownerName = GetRoleStr(
                            GameRoomM.Instance.getPlayerInfoByID(role.OwnerUid).Nickname,
                            role.OwnerUid == RoleM.Instance.Uid,
                            role.GuildId
                        );
                        displayName = $"{ownerName}\n{role.roleName}";
                    }
                    else
                    {
                        displayName = GetRoleStr($"{(string.IsNullOrEmpty(role.GuildName) ? string.Empty : "[" + role.GuildName + "]")}\n{role.roleName}", uid == RoleM.Instance.Uid, role.GuildId);
                    }

                    hud.text.text = displayName;

                    if (hud.img_title != null)
                        GameUtil.Instance.SetTitleIcon(role.TitleId, hud.img_title);
                }

                ProcessBuffers(hud, role._BufferList);
            }
        }
    }

    // 公共的buff处理逻辑
    private void ProcessBuffers(HUD hud, List<int> bufferList)
    {

        if (hud.Img_buffers == null) return;

        int bufferCount = bufferList?.Count ?? 0;
        int maxToProcess = Mathf.Min(hud.Img_buffers.Length, bufferCount);

        for (int j = 0; j < hud.Img_buffers.Length; j++)
        {
            var img = hud.Img_buffers[j];
            bool shouldShow = j < maxToProcess;
            img.gameObject.SetActive(shouldShow);

            if (shouldShow)
            {
                int buffId = bufferList[j];
                var cfgBuff = CFG_BuffEffect.GetById(buffId);

                if (!string.IsNullOrEmpty(cfgBuff.Icon))
                {
                    ResouseHelper.Instance.LoadSpriteByAtlas(
                        cfgBuff.Gallery,
                        cfgBuff.Icon,
                        img
                    );
                }
            }
        }
    }

    private string GetRoleStr(string str, bool isMy, long GuildId)
    {
        int QuaIdx = 8;
        if (GameRoomM.Instance.GameAttack == GameRoomM.AttackClass.Player)
        {
            if (isMy) //自己
                QuaIdx = 2;
            else //敌人
                QuaIdx = 7;
        }
        else if (GameRoomM.Instance.GameAttack == GameRoomM.AttackClass.Team)
        {
            if (isMy || (GuildId != 0 && GuildId == MainRole._MainRole.GuildId)) //自己
                QuaIdx = 2;
            else //敌人
                QuaIdx = 7;
        }

        return GameUtil.Instance.GetStrByQuality(QuaIdx, str);
    }

    private Camera mainCamera, uiCamera;

    private Vector2 GetCanvarPos(Vector3 worldPoint)
    {
        Vector2 anchoredPosition;

        if (mainCamera == null || uiCamera == null)
        {
            if (GameRoomM.Instance.IsTurnBasedRoom())
            {
                mainCamera = TurnBasedFightCtrl.Instance.MainCamer;
                uiCamera = TurnBasedFightCtrl.Instance.UICamera;
            }
            else
            {
                mainCamera = FightController.Instance.MainCamer;
                uiCamera = FightController.Instance.UICamera;
            }
        }

        Vector2 screenPosition = screenPosition = mainCamera.WorldToScreenPoint(worldPoint);
        RectTransformUtility.ScreenPointToLocalPointInRectangle(_UIViewGame3D.rect_Monster, screenPosition, uiCamera, out anchoredPosition);
        //anchoredPosition.y += 20;
        return anchoredPosition;
    }

    private struct HPEffItem
    {
        public Transform baseRoot;
        public Text txt_PlayerJian;
        public Text txt_Add;
        public Text txt_Enemy;
    }
    private Queue<HPEffItem> HpEffQueue = new Queue<HPEffItem>();

    /// <summary>
    /// 血量变化
    /// </summary>
    /// <param name="id">uid</param>
    /// <param name="changeNum">变化值</param>
    /// <param name="isTurnBased">是否为回合游戏</param>
    private void OnHPChangeEff(long id, int changeNum, bool isTurnBased = false)
    {
        if (HuDDic.ContainsKey(id))
        {
            HPEffItem hPEffItem;
            if (HpEffQueue.Count > 0)
            {
                hPEffItem = HpEffQueue.Dequeue();
            }
            else
            {
                GameObject go = GameObject.Instantiate(_UIViewGame3D.root_FlyHP.gameObject, _UIViewGame3D.rect_HPChnage);
                hPEffItem = new HPEffItem
                {
                    baseRoot = go.transform,
                    txt_PlayerJian = JUIExtensions.txt(go.transform),
                    txt_Add = JUIExtensions.txt(go.transform.Find("add")),
                    txt_Enemy = JUIExtensions.txt(go.transform.Find("jian_monster"))
                };
            }
            var basePos = HuDDic[id].obj.transform.localPosition;
            hPEffItem.baseRoot.transform.localPosition = new Vector3(basePos.x, basePos.y + 50, basePos.z);
            hPEffItem.baseRoot.gameObject.SetActive(true);
            bool isEnemy;
            if (isTurnBased)
                isEnemy = TurnBasedRoleManager.Instance.GetRoleCampType(id) == TurnBasedBattle.ActionCamp.RED_CAMP;
            else
                isEnemy = FightController.Instance.GetOneRoleInfo(id) == null || FightController.Instance.GetOneRoleInfo(id).role.isEnemy;

            if (!isEnemy && GameRoomM.Instance.GameAttack == GameRoomM.AttackClass.Player)
            {
                if (isTurnBased)
                {
                    var hasRlayer = TurnBasedRoleManager.Instance.GetRoleByUid(id, out var role);
                    if (hasRlayer)
                        isEnemy = role.campType == TurnBasedBattle.ActionCamp.RED_CAMP;
                }
            }

            hPEffItem.txt_PlayerJian.text = !isEnemy && changeNum < 0 ? changeNum.ToString() : "";
            hPEffItem.txt_Add.text = changeNum > 0 ? new StringBuilder("+").Append(changeNum).ToString() : "";
            hPEffItem.txt_Enemy.text = isEnemy && changeNum < 0 ? changeNum.ToString() : "";

            int off = changeNum > 0 ? -30 : 50;
            if (changeNum < 0)
            {
                int offx = off;
                if (UnityEngine.Random.Range(0, 2) == 0)
                    offx = -offx;
                DOTween.Sequence()
                    .Append(hPEffItem.baseRoot.DOScale(Vector3.one * 1.2f, 0.2f)).
                        Append(hPEffItem.baseRoot.DOLocalMoveY(hPEffItem.baseRoot.localPosition.y + off * 2, 0.3f)).
                        Append(hPEffItem.baseRoot.DOScale(Vector3.one * 1.5f, 0.2f)).
                        Append(hPEffItem.baseRoot.DOLocalMove(new Vector3(hPEffItem.baseRoot.localPosition.x + offx, hPEffItem.baseRoot.localPosition.y + (off * 3), hPEffItem.baseRoot.localPosition.z), 0.2f)).
                        OnComplete(() =>
                        {
                            hPEffItem.baseRoot.DOKill();
                            hPEffItem.baseRoot.gameObject.SetActive(false);
                            HpEffQueue.Enqueue(hPEffItem);
                        });
            }
            else
                hPEffItem.baseRoot.transform.DOLocalMoveY(hPEffItem.baseRoot.transform.localPosition.y + off, 0.3f).OnComplete(() =>
                {
                    hPEffItem.baseRoot.DOKill();
                    hPEffItem.baseRoot.gameObject.SetActive(false);
                    HpEffQueue.Enqueue(hPEffItem); ;
                });
        }
    }

    private async UniTaskVoid OnMissOrCriticalShow(long id, bool miss, bool critical)
    {
        if (HuDDic.ContainsKey(id) && (miss || critical))
        {
            Vector3 endPos = HuDDic[id].obj.transform.position;
            string poolKey = miss ? "Attack_Miss" : "Attack_BJ";
            if (!PoolHelper.Instance.CheckPool(poolKey))
            {
                var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
                await UniTask.WaitUntil(() => baseObj.Loaded);
                PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
            }
            await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
            GameObject go = PoolHelper.Instance.GetPoolObj(poolKey);

            go.transform.SetParent(_UIViewGame3D.rect_HPChnage);
            go.transform.localScale = Vector3.one;
            go.transform.position = endPos;
            go.gameObject.SetActive(true);

            if (miss)
                go.transform.DOLocalMoveY(go.transform.localPosition.y + (-50), 0.2f).OnComplete(() =>
                {
                    go.transform.DOKill();
                    PoolHelper.Instance.RecycObj(go);
                });
            else
            {
                DOTween.Sequence().Append(go.transform.DOScale(Vector3.one * 1.5f, 0.2f)).
                       Append(go.transform.DOLocalMoveY(go.transform.localPosition.y + 90, 0.3f)).
                       OnComplete(() =>
                       {
                           go.transform.DOKill();
                           PoolHelper.Instance.RecycObj(go);
                       });
            }
        }
    }

    #region 小地图

    private Dictionary<MapType, Queue<GameObject>> MapTipsQueueDic = new Dictionary<MapType, Queue<GameObject>>();
    private Dictionary<long, GameObject> MapTipsDic = new Dictionary<long, GameObject>();
    private enum MapType
    {
        TeamPos,
        EnemyPlayerPos,
        BossPos,
    }
    private GameObject GetMapObjectTip(MapType mapType)
    {
        GameObject go;
        if (MapTipsQueueDic.ContainsKey(mapType) && MapTipsQueueDic[mapType].Count > 0)
        {
            go = MapTipsQueueDic[mapType].Dequeue();
        }
        else
        {
            var obj = mapType switch
            {
                MapType.TeamPos => GetView().root_TeamPos.gameObject,
                MapType.EnemyPlayerPos => GetView().root_EnemyPlayerPos.gameObject,
                MapType.BossPos => GetView().root_BossPos.gameObject,
                _ => null
            };
            go = GameObject.Instantiate(obj, GetView().root_Map);
            go.gameObject.name = ((int)mapType).ToString();
        }
        go.gameObject.SetActive(true);
        return go;
    }

    private void RecycleMapObjectTip(long uid)
    {
        if (MapTipsDic.ContainsKey(uid) && uid != RoleM.Instance.Uid)
        {
            MapTipsDic[uid].SetActive(false);

            var mapType = (MapType)int.Parse(MapTipsDic[uid].name);
            if (!MapTipsQueueDic.ContainsKey(mapType))
                MapTipsQueueDic[mapType] = new Queue<GameObject>();

            MapTipsQueueDic[mapType].Enqueue(MapTipsDic[uid]);

            MapTipsDic.Remove(uid);
        }
    }

    private double _OneOff = 140 / 31f;
    private int _HalfMapDis = 16;

    private Vector3 GetMapTipPos(int x, int y)
    {
        Vector3 point = Vector3.zero;

        int MaxX = SceneManager.Wight - _HalfMapDis;
        if (MainRole._MainRole.pointNow.X >= MaxX)
        {
            point.x = (float)((x - MaxX) * _OneOff);
        }
        else if (MainRole._MainRole.pointNow.X > _HalfMapDis)
        {
            x -= (int)MainRole._MainRole.pointNow.X;
            point.x = (float)(x * _OneOff);
        }
        else
        {
            point.x = (float)((x - _HalfMapDis) * _OneOff);
        }

        int MaxY = SceneManager.Height - _HalfMapDis;
        if (MainRole._MainRole.pointNow.Y >= MaxY)
        {
            point.y = (float)((y - MaxY) * _OneOff);
        }
        else if (MainRole._MainRole.pointNow.Y > _HalfMapDis)
        {
            y -= (int)MainRole._MainRole.pointNow.Y;
            point.y = (float)(y * _OneOff);
        }
        else
        {
            point.y = (float)((y - _HalfMapDis) * _OneOff);
        }
        return point;
    }

    private Vector3 GetMapBGPos(int x, int y)
    {
        float _One = -470.0f / (SceneManager.Wight - 1);

        var pointNow = new Vector3(-60 + x * _One, -60 + y * _One, 0);
        if (pointNow.x > -120)
            pointNow.x = -120;
        else if (pointNow.x < -470)
            pointNow.x = -470;

        if (pointNow.y > -120)
            pointNow.y = -120;
        else if (pointNow.y < -470)
            pointNow.y = -470;
        return pointNow;
    }

    private double getModleDistance(double x, double y)
    {
        var pointNow = MainRole._MainRole.pointNow;
        return Math.Max(Math.Abs(x - pointNow.X), Math.Abs(y - pointNow.Y));
    }
    private void OnRefreshMap()
    {
        if (!_UIViewGame3D.root_Pos.activeSelf)
            return;

        var players = FightController.Instance.GetRolesInfo();
        int playerCount = players.Count;
        for (int i = 0; i < playerCount; i++)
        {
            var playerInfo = players[i];
            if (playerInfo?.role == null) continue;

            var role = playerInfo.role;
            if (role.uID != MainRole._MainRole.uID && (role._ObjectType != pb.ObjectType.PLAYER || getModleDistance(role._gridNow.X, role._gridNow.Y) > 16 || role.hp <= 0))
            {
                if (MapTipsDic.ContainsKey(role.uID))
                    RecycleMapObjectTip(role.uID);
                continue;
            }

            GameObject go = null;
            if (MapTipsDic.ContainsKey(role.uID))
                go = MapTipsDic[role.uID];
            else
            {
                go = role.uID == MainRole._MainRole.uID ? GetView().root_MyPos.gameObject :
                    role.isEnemy ? GetMapObjectTip(MapType.EnemyPlayerPos) : GetMapObjectTip(MapType.EnemyPlayerPos);
                MapTipsDic[role.uID] = go;
            }
            go.transform.localPosition = GetMapTipPos((int)role._gridNow.X, (int)role._gridNow.Y);

            if (role.uID == MainRole._MainRole.uID)
                GetView().root_MapBG.localPosition = GetMapBGPos((int)role._gridNow.X, (int)role._gridNow.Y);
        }

        var monsters = MonsterManager.Instance.getMonsterObjs();
        int monsterCount = monsters.Count;
        for (int i = 0; i < monsterCount; i++)
        {
            var monsterInfo = monsters[i];
            if (monsterInfo?._monser == null) continue;

            var monser = monsterInfo?._monser;
            if ((monser._MonsterCFG.Isboss != 1 && monser._MonsterCFG.Isboss != 4) || getModleDistance(monser._gridNow.X, monser._gridNow.Y) > 16)
            {
                if (MapTipsDic.ContainsKey(monser.uID))
                    RecycleMapObjectTip(monser.uID);
                continue;
            }

            GameObject go = null;
            if (MapTipsDic.ContainsKey(monser.uID))
            {
                go = MapTipsDic[monser.uID];
            }
            else
            {
                go = GetMapObjectTip(MapType.BossPos);
                MapTipsDic[monser.uID] = go;
            }
            go.transform.localPosition = GetMapTipPos((int)monser._gridNow.X, (int)monser._gridNow.Y);
        }
    }



    #endregion

    public override void OnExit()
    {
        _UIViewGame3D.UpateCallBack = null;
        HpEffQueue.Clear();
        CleatHUB();
        OnEventAddOrReMove(false);
    }

    private void OnEventAddOrReMove(bool isAdd)
    {
        if (isAdd)
        {
            UniGameEvent.AddEventListener<GameEvent.Fight_Die>(OnHUDFight_Die);
            UniGameEvent.AddEventListener<GameEvent.Fight_HPChange>(OnHPChangeEff);
            UniGameEvent.AddEventListener<GameEvent.Update_OneMonster>(OnRefresh_OneMonster);
        }
        else
        {
            UniGameEvent.RemoveEventListener<GameEvent.Fight_Die>(OnHUDFight_Die);
            UniGameEvent.RemoveEventListener<GameEvent.Fight_HPChange>(OnHPChangeEff);
            UniGameEvent.RemoveEventListener<GameEvent.Update_OneMonster>(OnRefresh_OneMonster);
        }
    }

    private void OnRefresh_OneMonster(GameEvent.Update_OneMonster _e)
    {
        if (HuDDic.ContainsKey(_e.uid))
        {
            HuDDic[_e.uid].IsUsed = false;
        }
    }


    private void OnHUDFight_Die(GameEvent.Fight_Die _e)
    {
        SetHUDById(_e.id);
        RecycEnemyTip(_e.id);
        RecycleMapObjectTip(_e.id);
    }

    private void OnHPChangeEff(GameEvent.Fight_HPChange _e)
    {
        //if (GameRoomM.Instance.CheakIsMainRole(_e.tid))
        //{
        //    Debug.Log($"_e.AttackerId  {_e.AttackerId} {!GameRoomM.Instance.CheakIsMainRole(_e.AttackerId)} {_e.changeNum} ");
        //}
        if (GameRoomM.Instance.GameAttack == GameRoomM.AttackClass.Team &&
            !GameRoomM.Instance.CheakIsMainRole(_e.tid) && !GameRoomM.Instance.CheakIsMainRole(_e.AttackerId))
            return;
        OnHPChangeEff(_e.tid, _e.changeNum, _e.isTurnBased);
        OnMissOrCriticalShow(_e.tid, _e.Miss, _e.Critical).Forget();
    }

    public override void OnDestroy()
    {
        OnExit();
    }
}