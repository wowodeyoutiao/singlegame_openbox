using Cysharp.Threading.Tasks;

using DG.Tweening;

using HotUpdateScripts.ProtoM;

using JEngine.Core;
using JEngine.Event;
using JEngine.Helper;

using MetaFramework.UI;

using pb;

using System;
using System.Collections.Generic;

using UnityEngine;

using Object = UnityEngine.Object;

namespace HotUpdateScripts.Game
{
    public class TurnBasedFightCtrl : MonoBehaviour
    {
        public enum BattleState
        {
            None,
            Init,
            SelectAttacker,
            Attack,
            AttackCompleted,
            Over,
        }

        // 战斗状态
        public BattleState currentBattleState;

        // 当前攻击者角色
        public TurnBasedRole currentBattleRole;

        private readonly string _positionRootPath = "Common/Prefab/";

        public int battleWaitTime = 3;

        private Vector3 m_BlueRootPos_PVP => new Vector3(-4, -2.5f, -2);
        private Vector3 m_RedRootPos_PVP => new Vector3(4, -2.5f, -2);

        private Vector3 m_BlueRootPos_WuXing => new Vector3(-4, -2.5f, -2);
        private Vector3 m_RedRootPos_WuXing => new Vector3(4, -2.5f, -2);


        private static TurnBasedFightCtrl _instance;
        public static TurnBasedFightCtrl Instance
        {
            get
            {
                if (_instance == null)
                {
                    GameObject go = new GameObject("FightCtrl");
                    _instance = go.AddComponent<TurnBasedFightCtrl>();
                }
                return _instance;
            }
        }

        private Camera _MainCamer;
        public Camera MainCamer
        {
            get
            {
                if (_MainCamer == null)
                {
                    _MainCamer = GameObject.Find("Main Camera").GetComponent<Camera>();
                }
                return _MainCamer;
            }
        }

        private Camera _UICamera;
        public Camera UICamera
        {
            get
            {
                if (_UICamera == null)
                {
                    _UICamera = GameObject.Find("UICamera").GetComponent<Camera>();
                }
                return _UICamera;
            }
        }

        //特效节点
        private Transform _effectRoot;
        public Transform EffectRoot
        {
            get
            {
                if (_effectRoot == null)
                {
                    _effectRoot = new GameObject("EffectRoot").transform;
                    _effectRoot.position = new Vector3(0, 0, -3);
                }
                return _effectRoot;
            }
        }

        private Transform _roleRoot;
        public Transform RoleRoot
        {
            get
            {
                if (_roleRoot == null)
                {
                    GameObject roleRoot = GameObject.Find("RoleRoot");
                    if (!roleRoot)
                    {
                        roleRoot = new GameObject("RoleRoot");
                    }
                    _roleRoot = roleRoot.transform;
                }
                return _roleRoot;
            }
        }

        //蓝方人物
        private Transform _blueRoleRoot;
        public Transform BlueRoleRoot
        {
            get
            {
                if (_blueRoleRoot == null)
                {
                    _blueRoleRoot = RoleRoot.Find("BlueRoot(Clone)");
                }
                return _blueRoleRoot;
            }
        }

        //红方人物
        private Transform _redRoleRoot;
        public Transform RedRoleRoot
        {
            get
            {
                if (_redRoleRoot == null)
                {
                    _redRoleRoot = RoleRoot.Find("RedRoot(Clone)");
                }
                return _redRoleRoot;
            }
        }

        private SpriteRenderer mapRender;

        public float attackMoveSpeed = 0.15f;

        public TurnBasedBattle.TBPlayer bluePlayer;

        public TurnBasedBattle.TBPlayer redPlayer;

        TimeItem m_StartBattleTimeItem;

        public void ClearStartBattleTime()
        {
            if (GameHelper.Instance != null && m_StartBattleTimeItem != null)
            {
                GameHelper.Instance.RemoveInvoke(m_StartBattleTimeItem);
                m_StartBattleTimeItem = null;
            }
        }

        public async UniTask OnInit()
        {
            GameRoomM.Instance._GameState = GameRoomM.GameState.Ready;
            GameRoomM.Instance.CFGRoom = CFG_Room.GetById(TurnBasedM.Instance.roomId);

            AudioManager.Instance.PlayBgMusic((AudioKey.Sound)Enum.Parse(typeof(AudioKey.Sound), GameRoomM.Instance.CFGRoom.Music));
            bluePlayer = TurnBasedM.Instance.GetPlayerByCampType(TurnBasedBattle.ActionCamp.BLUE_CAMP);
            redPlayer = TurnBasedM.Instance.GetPlayerByCampType(TurnBasedBattle.ActionCamp.RED_CAMP);
            MainCamer.orthographicSize = 22;

            await LoadMap();

            UITool.Open<UIGame3D>();
            UITool.Open<UIGameTurnBased>((ui) =>
            {
                UITool.Open<UIPVPBeginTime>(null, TurnBasedM.Instance.GetFirstAttackCamp() == TurnBasedBattle.ActionCamp.BLUE_CAMP ?
                    UIPVPBeginTime.ShowType.MyTeamFirstAttack : UIPVPBeginTime.ShowType.OtherTeamFirstAttack);
            });

            await LoadPositionRoot();
            UITool.Close<UILoading>();

            MainCamer.DOOrthoSize(20, 1f).SetEase(Ease.OutQuad);

            currentBattleState = BattleState.None;

            m_StartBattleTimeItem = GameHelper.Instance.AddInvoke(new TimeItem(battleWaitTime, () =>
            {
                if (Object.Equals(TurnBasedRoleManager.Instance, null))
                    return;

                ChangeState(BattleState.Init);
                JEvent.defaultEvent.Post(new GameEvent.PVP_JumpGameCountdown());
            }));
        }

        public async UniTask LoadMap()
        {
            mapRender = GameObject.Find("Map").GetComponent<SpriteRenderer>();
            mapRender.transform.localPosition = new Vector3(0, -3, 0);
            CFGItem_Map MapCfg = CFG_Map.GetById(GameRoomM.Instance.CFGRoom.MapId);
            await ResouseHelper.Instance.LoadSprite($"Terrain/Texture/Map_{MapCfg.MapID}", mapRender);
        }

        public async UniTask LoadPositionRoot()
        {
            await ResouseHelper.Instance.LoadGameObj(_positionRootPath + "BlueRoot", o =>
            {
                _blueRoleRoot = Instantiate(o, RoleRoot).transform;

                if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVP_More)
                    _blueRoleRoot.localPosition = m_BlueRootPos_PVP;
                else if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.WuXing)
                    _blueRoleRoot.localPosition = m_BlueRootPos_WuXing;
            });
            await ResouseHelper.Instance.LoadGameObj(_positionRootPath + "RedRoot", o =>
            {
                _redRoleRoot = Instantiate(o, RoleRoot).transform;

                if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVP_More)
                    _redRoleRoot.localPosition = m_RedRootPos_PVP;
                else if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.WuXing)
                    _redRoleRoot.localPosition = m_RedRootPos_WuXing;
            });

            await UniTask.WaitUntil(() => _blueRoleRoot != null && _redRoleRoot != null);

            await LoadRole(TurnBasedBattle.ActionCamp.BLUE_CAMP);
            await LoadRole(TurnBasedBattle.ActionCamp.RED_CAMP);
        }

        public async UniTask LoadRole(TurnBasedBattle.ActionCamp actionCamp)
        {
            var players = actionCamp switch
            {
                TurnBasedBattle.ActionCamp.BLUE_CAMP => bluePlayer,
                TurnBasedBattle.ActionCamp.RED_CAMP => redPlayer,
                _ => null,
            };

            for (int i = 0; i < players.Objs.Count; i++)
            {
                if (players.Objs[i].Uid == 0 || players.Objs[i].Type == ObjectType.INVALID_ObjectType)
                {
                    continue;
                }

                string poolKey = "";
                string poolPath = "";

                switch (players.Objs[i].Type)
                {
                    case ObjectType.PLAYER:
                        poolKey = "MainRoleFight_" + players.Objs[i].Job + players.Objs[i].Gender;
                        poolPath = GameConst.GetModelByJobAndSex((JobType)players.Objs[i].Job, (SexType)players.Objs[i].Gender, true);
                        break;
                    case ObjectType.PARTNER:
                        poolKey = "Fight_Partner_" + players.Objs[i].CfgID;
                        poolPath = GameConst.GetMonsterModelByID(poolKey);
                        break;
                    case ObjectType.MONSTER:
                        poolKey = CFG_Monster.GetById(players.Objs[i].CfgID).Model;
                        poolPath = GameConst.GetMonsterModelByID(poolKey);
                        break;
                    default:
                        break;
                }

                if (!PoolHelper.Instance.CheckPool(poolKey))
                {
                    var baseObj = new JPrefab(poolPath, true);
                    await UniTask.WaitUntil(() => baseObj.Loaded);
                    PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
                }
                await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
                GameObject roleModelObj = PoolHelper.Instance.GetPoolObj(poolKey);
                roleModelObj.gameObject.name = poolKey + "_" + players.Objs[i].Uid;

                TurnBasedRole turnBasedRole = roleModelObj.GetComponent<TurnBasedRole>();
                if (!turnBasedRole) turnBasedRole = roleModelObj.AddComponent<TurnBasedRole>();

                int position = players.Objs[i].Position + 1;

                switch (actionCamp)
                {
                    case TurnBasedBattle.ActionCamp.BLUE_CAMP:
                        roleModelObj.transform.SetParent(BlueRoleRoot.Find($"Pos_{position}"));
                        roleModelObj.transform.localRotation = Quaternion.Euler(0, 180, 0);
                        turnBasedRole.campType = TurnBasedBattle.ActionCamp.BLUE_CAMP;
                        TurnBasedRoleManager.Instance.AddBuleRole(turnBasedRole, position);

                        break;
                    case TurnBasedBattle.ActionCamp.RED_CAMP:
                        roleModelObj.transform.SetParent(RedRoleRoot.Find($"Pos_{position}"));
                        roleModelObj.transform.localRotation = Quaternion.Euler(0, 0, 0);
                        turnBasedRole.campType = TurnBasedBattle.ActionCamp.RED_CAMP;
                        TurnBasedRoleManager.Instance.AddRedRole(turnBasedRole, position);

                        break;
                    default:
                        break;
                }

                roleModelObj.transform.localPosition = Vector3.zero;
                turnBasedRole.uid = players.Objs[i].Uid;
                turnBasedRole.cfgId = players.Objs[i].CfgID;
                turnBasedRole.roleName = players.Objs[i].Name;
                turnBasedRole.maxHp = players.Objs[i].MaxHp;
                turnBasedRole.currentHp = players.Objs[i].CurHp;
                turnBasedRole.job = players.Objs[i].Job;
                turnBasedRole.position = position;
                turnBasedRole.roleType = players.Objs[i].Type;
                turnBasedRole.mainElement = players.Objs[i].ElementMain;

                switch (players.Objs[i].Type)
                {
                    case ObjectType.PLAYER:
                        roleModelObj.transform.localScale = Vector3.one * 0.8f;
                        break;
                    case ObjectType.PARTNER:
                    case ObjectType.MONSTER:
                        roleModelObj.transform.localScale = Vector3.one * GameConst.GetMonsterScalModelByID(poolKey);
                        break;
                    default:
                        break;
                }

                roleModelObj.SetActive(true);
            }

            JEvent.defaultEvent.Post(new GameEvent.PVP_TurnBasedShowUserName { isMy = actionCamp == TurnBasedBattle.ActionCamp.BLUE_CAMP });
        }

        /// <summary>
        /// 战斗状态
        /// </summary>
        /// <param name="battleState"></param>
        public void ChangeState(BattleState battleState)
        {
            currentBattleState = battleState;

            switch (battleState)
            {
                case BattleState.None:
                    break;
                case BattleState.Init:
                    StartGame();
                    break;
                case BattleState.SelectAttacker:
                    SelectAttacker();
                    break;
                case BattleState.Attack:
                    OnAttack();
                    break;
                case BattleState.AttackCompleted:
                    AttackCompleted();
                    break;
                case BattleState.Over:
                    OnPopupResult();
                    break;
                default:
                    break;
            }
        }

        // 总回合数
        int m_TotalRound;
        // 当前回合
        int m_CurrentRound;
        // 总的攻击战斗次数
        int m_AllAttackNum;
        // 当前第几次战斗
        int m_CurrentAttackNum;
        // 是否需要延迟攻击？
        bool m_NeedDelayAttack;

        public int totalRound => m_TotalRound;
        public int currentRound => m_CurrentRound;
        public int attackNum => m_CurrentAttackNum;

        TurnBasedBattle.Round m_CurrendRoundData;
        //TurnBasedBattle.Round m_LastRoundData;

        // 开始回合
        public void StartGame()
        {
            m_TotalRound = TurnBasedM.Instance.turnBasedBattleData.Rounds[TurnBasedM.Instance.turnBasedBattleData.Rounds.Count - 1].RoundId;
            m_CurrentRound = 1;
            m_AllAttackNum = TurnBasedM.Instance.turnBasedBattleData.Rounds.Count;
            m_CurrentAttackNum = 1;
            ChangeState(BattleState.SelectAttacker);
        }

        public bool IsBattleEnd()
        {
            return m_CurrentAttackNum >= m_AllAttackNum;
        }

        TimeItem delayAttackTime;

        public void ClearDelayAttackTime()
        {
            if (GameHelper.Instance != null && delayAttackTime != null)
            {
                GameHelper.Instance.RemoveInvoke(delayAttackTime);
                delayAttackTime = null;
            }
        }

        public void AttackCompleted()
        {
            if (IsBattleEnd())
            {
                ChangeState(BattleState.Over);
            }
            else
            {
                m_CurrentAttackNum++;
                if (m_NeedDelayAttack)
                {
                    ClearDelayAttackTime();
                    delayAttackTime = GameHelper.Instance.AddInvoke(new TimeItem(1, () =>
                    {
                        m_NeedDelayAttack = false;
                        ChangeState(BattleState.SelectAttacker);
                        ClearDelayAttackTime();
                    }));
                }
                else
                    ChangeState(BattleState.SelectAttacker);
            }
        }

        public void SelectAttacker()
        {
            m_CurrendRoundData = GetCurrentRoundData();
            m_CurrentRound = m_CurrendRoundData.RoundId;
            //m_LastRoundData = GetLastRoundData();

            switch (m_CurrendRoundData.ActionCamp)
            {
                case TurnBasedBattle.ActionCamp.BLUE_CAMP:
                    if (TurnBasedRoleManager.Instance.GetBlueRoleByUid(m_CurrendRoundData.Uid, out TurnBasedRole blueRole))
                        currentBattleRole = blueRole;

                    break;
                case TurnBasedBattle.ActionCamp.RED_CAMP:
                    if (TurnBasedRoleManager.Instance.GetRedRoleByUid(m_CurrendRoundData.Uid, out TurnBasedRole redRole))
                        currentBattleRole = redRole;

                    break;
                default:
                    break;
            }

            if (currentBattleRole == null)
            {
                return;
            }

            currentBattleRole.isDead = false;
            currentBattleRole.skillId = m_CurrendRoundData.SkillId;
            ChangeState(BattleState.Attack);
        }

        private void OnAttack()
        {
            JEvent.defaultEvent.Post(new GameEvent.PVP_TurnBasedRoundShow
            {
                currentRound = m_CurrentRound,
                totalRound = m_TotalRound,
            });

            currentBattleRole.StartAttack();
        }

        public TurnBasedBattle.Round GetCurrentRoundData()
        {
            return TurnBasedM.Instance.GetAttackAction(m_CurrentAttackNum);
        }

        public TurnBasedBattle.Round GetLastRoundData()
        {
            if (m_CurrentAttackNum == 1)
                return null;

            return TurnBasedM.Instance.GetAttackAction(m_CurrentAttackNum - 1);
        }

        // 添加buff
        public void AddBuffers()
        {
            bool hasRevive = false;
            List<long> uids = new List<long>();

            foreach (var effect in m_CurrendRoundData.Effects)
            {
                if (effect.TargetUid != 0)
                {
                    if (TurnBasedRoleManager.Instance.GetRoleByUid(effect.TargetUid, out var role))
                    {
                        if (effect.Addbuffs != null)
                            role.AddBuffers(effect.Addbuffs);

                        // 有回血的buff时，先更新血量血条
                        if (IsHeal(effect.Addbuffs))
                        {
                            var objData = GetCurrentRoleData(effect.TargetUid);
                            if (objData != null)
                            {
                                role.currentHp = objData.CurHp;
                                role.maxHp = objData.MaxHp;
                            }

                            UniGameEvent.Send(new GameEvent.Fight_HPChange
                            {
                                tid = effect.TargetUid,
                                changeNum = effect.Damage > 0 ? -effect.Damage : effect.Recover,
                                isTurnBased = true,
                            });
                        }

                        // 复活时忽略先死亡再立即复活的角色，这些角色的飘字和动作推迟到被攻击时表现。
                        if (effect.Revive && !IsDeadAtCurrentAttack(effect.TargetUid))
                        {
                            hasRevive = true;
                            m_NeedDelayAttack = true;
                            role.Revive().Forget();
                            uids.Add(effect.TargetUid);

                            //UniGameEvent.Send(new GameEvent.Fight_HPChange
                            //{
                            //    tid = effect.TargetUid,
                            //    changeNum = effect.Damage > 0 ? -effect.Damage : effect.Recover,
                            //    Critical = effect.Crit,
                            //    Miss = effect.Miss,
                            //    isTurnBased = true,
                            //});
                        }

                    }
                }
            }

            // 如果有复活的角色，更新角色血量
            if (hasRevive)
            {
                foreach (var uid in uids)
                {
                    if (TurnBasedRoleManager.Instance.GetRoleByUid(uid, out var _role))
                    {
                        var objData = GetCurrentRoleData(uid);
                        if (objData != null)
                        {
                            _role.currentHp = objData.CurHp;
                            _role.maxHp = objData.MaxHp;
                        }
                    }
                }
            }
        }

        public bool HasRevive()
        {
            foreach (var effect in m_CurrendRoundData.Effects)
            {
                if (effect.TargetUid != 0 && effect.Revive)
                {
                    return true;
                }
            }
            return false;
        }

        // 是否带有回血buff
        private bool IsHeal(int[] buffs)
        {
            if (buffs == null) return false;

            foreach (var buff in buffs)
            {
                // 一次性回血buff
                if (buff == 8 || buff == 22) return true;
            }
            return false;
        }

        public bool IsDeadAtCurrentAttack(long targetUid)
        {
            return m_CurrendRoundData.Effects.Exists(v => v.Dead && v.TargetUid == targetUid);
        }

        public void UpdateRoleInfo()
        {
            for (int i = 0; i < m_CurrendRoundData.RedStats.Count; i++)
            {
                var objData = m_CurrendRoundData.RedStats[i];
                if (objData.Uid != 0 && TurnBasedRoleManager.Instance.GetRedRoleByUid(objData.Uid, out var role))
                {
                    role.maxHp = objData.MaxHp;

                    if (IsDeadAtCurrentAttack(objData.Uid))
                        role.currentHp = 0;
                    else
                        role.currentHp = objData.CurHp;
                }
            }

            for (int i = 0; i < m_CurrendRoundData.BlueStats.Count; i++)
            {
                var objData = m_CurrendRoundData.BlueStats[i];
                if (objData.Uid != 0 && TurnBasedRoleManager.Instance.GetBlueRoleByUid(objData.Uid, out var role))
                {
                    role.maxHp = objData.MaxHp;

                    if (IsDeadAtCurrentAttack(objData.Uid))
                        role.currentHp = 0;
                    else
                        role.currentHp = objData.CurHp;
                }
            }

            foreach (var effect in m_CurrendRoundData.Effects)
            {
                // 被复活的角色提前飘字，不会等到攻击时飘字。
                if (!effect.Revive && !IsHeal(effect.Addbuffs))
                {
                    UniGameEvent.Send(new GameEvent.Fight_HPChange
                    {
                        tid = effect.TargetUid,
                        changeNum = effect.Damage > 0 ? -effect.Damage : effect.Recover,
                        Critical = effect.Crit,
                        Miss = effect.Miss,
                        isTurnBased = true,
                    });

                    Debug.Log($"effect.Crit = " + effect.Crit);
                }
            }
        }

        public TurnBasedBattle.TBObject GetCurrentRoleData(long uid)
        {
            var data = m_CurrendRoundData.BlueStats.Find(v => v.Uid == uid);
            data ??= m_CurrendRoundData.RedStats.Find(v => v.Uid == uid);
            return data;
        }

        public void ClearDelayReviveTime()
        {
            if (GameHelper.Instance != null && delayReviveTime != null)
            {
                GameHelper.Instance.RemoveInvoke(delayReviveTime);
                delayReviveTime = null;
            }
        }

        TimeItem delayReviveTime;

        public void CollectKilledRole()
        {
            List<long> killedUids = new List<long>();

            foreach (var item in m_CurrendRoundData.Effects)
            {
                if (item.Dead)
                {
                    killedUids.Add(item.TargetUid);

                    if (killedUids.Count > 0)
                    {
                        KilledRole(killedUids);
                    }

                    // 死亡后在本次出手的周期内立马被复活的角色，需要严迟攻击（等复活站起来后）。
                    bool reviveImmediately = IsReviveAtCurrentAttack(item.TargetUid);
                    if (reviveImmediately)
                    {
                        if (TurnBasedRoleManager.Instance.GetRoleByUid(item.TargetUid, out TurnBasedRole role))
                        {
                            m_NeedDelayAttack = true;
                            ClearDelayReviveTime();
                            delayReviveTime = GameHelper.Instance.AddInvoke(new TimeItem(1.2f, () =>
                            {
                                var data = GetCurrentRoleData(role.uid);
                                role.Revive().Forget();
                                role.currentHp = data.CurHp;
                                role.maxHp = data.MaxHp;

                                //UniGameEvent.Send(new GameEvent.Fight_HPChange
                                //{
                                //    tid = item.TargetUid,
                                //    changeNum = item.Damage > 0 ? -item.Damage : item.Recover,
                                //    isTurnBased = true,
                                //});
                                ClearDelayReviveTime();
                            }));
                        }
                    }
                }
            }
        }

        private void KilledRole(List<long> uids)
        {
            for (int i = 0; i < uids.Count; i++)
            {
                if (TurnBasedRoleManager.Instance.GetRoleByUid(uids[i], out TurnBasedRole role))
                {
                    role.Dead();
                    UniGameEvent.Send(new GameEvent.Fight_Die
                    {
                        id = role.uid
                    });
                }
            }
        }

        /// <summary>
        /// 本次出手被打死后立马复活
        /// </summary>
        /// <param name="effects">效果列表</param>
        /// <param name="targetUid">目标角色</param>
        /// <returns></returns>
        public bool IsReviveAtCurrentAttack(long targetUid)
        {
            return m_CurrendRoundData.Effects.Exists(v => v.Revive && v.TargetUid == targetUid);
        }

        public TurnBasedBattle.Round.Effect GetReviveEffect(long targetUid)
        {
            return m_CurrendRoundData.Effects.Find(v => v.Revive && v.TargetUid == targetUid);
        }

        public bool IsNeedPlayHurtEffect(int skillId)
        {
            var cfgSkill = CFG_Skill.GetById(skillId);
            if (cfgSkill == null) return false;
            if (skillId == 999900) return true;

            int skillEffId = cfgSkill.SkillID;
            return skillEffId == 1 || skillEffId == 14 || skillEffId == 27 ||
                skillEffId == 20000 || skillEffId == 20001 || skillEffId == 20002;
        }

        public void PlayHurt(long uid, int skillId)
        {
            if (TurnBasedRoleManager.Instance.GetRoleByUid(uid, out TurnBasedRole role))
            {
                role.PlayHurt(skillId).Forget();
            }
        }

        /// <summary>
        /// 战斗结果
        /// </summary>
        public void OnPopupResult()
        {
            var result = TurnBasedM.Instance.GetFightResult();
            List<ItemInfo> awards = new List<ItemInfo>();
            List<ItemInfo> otherAwards = null;

            if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.PVP_More)
            {
                awards = result.Items;
            }
            else if (GameRoomM.Instance.CFGRoom.Roomtype == CFG_Room.RoomType.WuXing)
            {
                if (result.Win)
                {
                    awards.Add(result.Items[0]);
                    awards.Add(result.Items[1]);

                    otherAwards = new List<ItemInfo>();
                    for (int i = 2; i < result.Items.Count; i++)
                    {
                        otherAwards.Add(result.Items[i]);
                    }
                }
            }

            var gameResult = new GameResultBattlePush
            {
                RoomId = TurnBasedM.Instance.roomId,
                BattleId = 0,
                OverType = result.Win ? GameOverType.GameOverTypeMDead : GameOverType.GameOverTypeDeadFaild,
                Items = awards
            };
            UITool.Open<UIGameFightOver>(null, gameResult, otherAwards);
        }

        /// <summary>
        /// 跳过战斗，直接显示结果。
        /// </summary>
        public void JumpGame()
        {
            TurnBasedRoleManager.Instance.StopAllRoleAttack();
            ChangeState(BattleState.Over);
        }

        // 回合赛：重新挑战，打开布阵界面。
        public void GoNextFightPVPMore()
        {
            OnExitGame(() =>
            {
                UITool.Open<UIFightChoice>(null, UIFightChoice.SelectType.Limit, UIFightChoice.ChanllengeType.PvP);
                UITool.Open<UIFightFormation>(null, GameRoomM.Instance.CFGRoom.Roomtype, TurnBasedM.Instance.opponentUid);
            });
        }

        // 下一关五行试炼关卡
        public void GoNextFightWuXing()
        {
            OnExitGame(() =>
            {
                UITool.Open<UIWuXingTrials>(null, new UIWuXingTrials.Params
                {
                    isPlayMoveLayerAnimation = true,
                    isContinueChallenge = true,
                });
            });
        }

        // 五行试炼赋能
        public void GoWuXingFuNeng()
        {
            OnExitGame(() =>
            {
                UITool.Open<UIWuXingTrials>(null, new UIWuXingTrials.Params
                {
                    isOpenFuNeng = true,
                });
            });
        }

        public void OnExitGame(Action JumpUI = null)
        {
            TurnBasedRoleManager.Instance.ClearAllRole();
            GameRoomM.Instance._GameState = GameRoomM.GameState.No;
            UITool.Open<UILoading>(null, new UILoading.Loading { _time = 0 });

            SenceChange handle = new SenceChange
            {
                scencePath = GameConst.HallPath,
                CallBack = () =>
                    {
                        UITool.Close<UILoading>();
                        UITool.Open<UIHall>(null);
                        JumpUI?.Invoke();
                    }
            };
            JEvent.defaultEvent.Post(handle);
        }

        void OnDestroy()
        {
            ClearStartBattleTime();
            ClearDelayReviveTime();
            ClearDelayAttackTime();
            TurnBasedRoleManager.Instance.OnDestroy();
        }
    }
}