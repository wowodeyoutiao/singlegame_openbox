using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using MetaFramework.UI;
using HotUpdateScripts.ProtoM;
using HotUpdateScripts.Game;
using pb;
using Spine;
using JEngine.UI;
using System;
using HotUpdateScripts;
using JEngine.Helper;

public class UIGameFightWarMap : UIBase<UIModelGameFightWarMap, UIViewGameFightWarMap>
{
    public override UILayer GetLayer()
    {
        return UILayer.Normal;
    }

    private float _OneX = 0;
    private float _OneY = 0;

    private float _BaseX = 0;
    private float _BaseY = 0;
    public override void OnInit()
    {
        GetView().btn_Close.onClick.AddListener(() =>
        {
            AudioManager.Instance.PlaySound(AudioKey.Sound.Button);
            UITool.Close<UIGameFightWarMap>();
        });

        GetView().txt_Tiltle.text = GameRoomM.Instance.CFGRoom.Name;

        GetView().touch_Map.AddPointerDownListener((eventData) =>
        {
            Vector2 localPoint;
            RectTransformUtility.ScreenPointToLocalPointInRectangle(GetView().rect_Map, eventData.position, JUIExtensions.uiCamera, out localPoint);
            MainRole._MainRole.targetCollectID = 0;
            MainRole._MainRole.targetID = 0;

            MainRole._MainRole._SendMovePos = new Point { X = (int)((localPoint.x - _BaseX) / _OneX), Y = (int)((localPoint.y - _BaseY) / _OneY) };
        });

        var CFGMap = CFG_Map.GetById(GameRoomM.Instance.CFGRoom.MapId);
        _OneX = GetView().rect_Map.rect.width / CFGMap.Wide;
        _OneY = GetView().rect_Map.rect.height / CFGMap.High;

        _BaseX = -GetView().rect_Map.rect.width / 2;
        _BaseY = -GetView().rect_Map.rect.height / 2;

        ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.FightAtlas, GameRoomM.Instance.CFGRoom.Name.Contains("地狱") ? "bg_qxzc_map_2" : "bg_qxzc_map_1", GetView().img_Map);
    }

    public override void OnLayerChange()
    {
        GetView().line.sortingOrder = _canvasOrder +1 ;
        GetView().line.startWidth = 0.03f;
        GetView().line.endWidth = 0.03f;
        GetView().line.positionCount = 2;
    }

    public override void OnEnter(params object[] args)
    {
        OnEventAddOrReMove(true);
        GetView().UpateCallBack = OnUpdate;
    }
    private enum MapType 
    {
        TeamPos,
        EnemyPlayerPos,
        BossPos,
    }

    private Dictionary<MapType,Queue<GameObject>> MapTipsQueueDic = new Dictionary<MapType, Queue<GameObject>>();
    private Dictionary<long, GameObject> MapTipsDic = new Dictionary<long, GameObject>();

    private GameObject GetGameObjectTip(MapType mapType)
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
                MapType.TeamPos => GetView().obj_MyTeam,
                MapType.EnemyPlayerPos => GetView().obj_PlayerEnemy,
                MapType.BossPos => GetView().obj_Boss,
                _ => null
            };
            go = GameObject.Instantiate(obj, GetView().rect_Map);
            go.gameObject.name = ((int)mapType).ToString();
        }
        go.gameObject.SetActive(true);
        return go;
    }


    private Vector3 GetPos(int x, int y)
    {
        return new Vector3(_BaseX + x * _OneX, _BaseY + y * _OneY, 0);
    }

    private void OnUpdate()
    {
        var players = FightController.Instance.GetRolesInfo();
        int playerCount = players.Count;
        for (int i = 0; i < playerCount; i++)
        {
            var playerInfo = players[i];
            if (playerInfo?.role == null) continue;

            var role = playerInfo.role;
            if (role._ObjectType != pb.ObjectType.PLAYER) continue;

            GameObject go = null;
            if (MapTipsDic.ContainsKey(role.uID))
            {
                go = MapTipsDic[role.uID];
            }
            else
            {
                go = role.uID == MainRole._MainRole.uID ? GetView().root_MyPos.gameObject :
                    role.isEnemy ? GetGameObjectTip(MapType.EnemyPlayerPos) : GetGameObjectTip(MapType.TeamPos);
                MapTipsDic[role.uID] = go;
            }
            go.transform.localPosition = GetPos((int)role._gridNow.X, (int)role._gridNow.Y);
        }

        var monsters = MonsterManager.Instance.getMonsterObjs();
        int monsterCount = monsters.Count;
        for (int i = 0; i < monsterCount; i++)
        {
            var monsterInfo = monsters[i];
            if (monsterInfo?._monser == null) continue;

            var monser = monsterInfo?._monser;
            if (monser._MonsterCFG.Isboss != 1 && monser._MonsterCFG.Isboss != 4) continue;

            GameObject go = null;
            if (MapTipsDic.ContainsKey(monser.uID))
            {
                go = MapTipsDic[monser.uID];
            }
            else
            {
                go = GetGameObjectTip(MapType.BossPos);
                MapTipsDic[monser.uID] = go;
            }
            go.transform.localPosition = GetPos((int)monser._gridNow.X, (int)monser._gridNow.Y);
        }

        GetView().root_ToPos.gameObject.SetActive(MainRole._MainRole._SendMovePos != default);
        GetView().line.gameObject.SetActive(MainRole._MainRole._SendMovePos != default);

        if (MainRole._MainRole._SendMovePos != default)
        {
            GetView().root_ToPos.localPosition = GetPos((int)MainRole._MainRole._SendMovePos.X, (int)MainRole._MainRole._SendMovePos.Y);
            GetView().txt_ToPos.text = $"({MainRole._MainRole._SendMovePos.X},{MainRole._MainRole._SendMovePos.Y})";
            GetView().line.SetPosition(0, MapTipsDic[MainRole._MainRole.uID].transform.position);
            GetView().line.SetPosition(1, GetView().root_ToPos.position);
        }
    }
    private void OnEventAddOrReMove(bool isAdd)
    {
        if (isAdd)
        {
            UniGameEvent.AddEventListener<GameEvent.Fight_Die>(OnHUDFight_Die);
        }
        else
        {
            UniGameEvent.RemoveEventListener<GameEvent.Fight_Die>(OnHUDFight_Die);
        }
    }

    private void RecycEnemyTip(long uid)
    {
        if (MapTipsDic.ContainsKey(uid)&& uid != RoleM.Instance.Uid)
        {
            MapTipsDic[uid].SetActive(false);

            var mapType = (MapType)int.Parse(MapTipsDic[uid].name);
            if (!MapTipsQueueDic.ContainsKey(mapType))
                MapTipsQueueDic[mapType] = new Queue<GameObject>();

            MapTipsQueueDic[mapType].Enqueue(MapTipsDic[uid]);

            MapTipsDic.Remove(uid);
        }
    }
    private void OnHUDFight_Die(GameEvent.Fight_Die _e)
    {
        RecycEnemyTip(_e.id);
    }

    public override void OnExit()
    {
        foreach (var item in MapTipsDic)
        {
            if (item.Key == MainRole._MainRole.uID) continue;
            item.Value.SetActive(false);
            var mapType = (MapType)int.Parse(item.Value.name);
            if (!MapTipsQueueDic.ContainsKey(mapType))
                MapTipsQueueDic[mapType] = new Queue<GameObject>();
            MapTipsQueueDic[mapType].Enqueue(item.Value);
        }
        MapTipsDic.Clear();
        GetView().UpateCallBack = null;
        OnEventAddOrReMove(false);
    }

    public override void OnDestroy()
    {
        GetView().UpateCallBack = null;
        OnEventAddOrReMove(false);
    }
}