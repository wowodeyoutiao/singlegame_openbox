using HotUpdateScripts;
using HotUpdateScripts.Game;
using HotUpdateScripts.ProtoM;
using JEngine.Helper;
using pb;
using System;
using System.Collections.Generic;
using UnityEngine;

public class MainRole : Role
{
    public static MainRole _MainRole { get { return _mainRole; } }
    private static MainRole _mainRole;

    public Direction _joyDir;

    /// <summary>
    /// 上次移动发送时间
    /// </summary>
    private float _sendMoveLast;
    /// <summary>
    /// 上次攻击发送时间
    /// </summary>
    private float _sendAttackLast;

    /// <summary>
    /// 攻击时间间隔
    /// </summary>
    private float _attackSpeed;

    private Point _lastSendLocation;
    private Point _lastMoveLocation;

    private long _targetID;//目标攻击玩家ID
    public long targetID
    {
        get => _targetID;

        set
        {
            if (_targetID == value)
                return;
            _targetID = value;

        }
    }

    private long _targetCollectID;//目标收集玩家ID
    public long targetCollectID
    {
        get => _targetCollectID;

        set
        {
            if (_targetCollectID == value)
                return;
            _targetCollectID = value;

        }
    }

    public Point _SendMovePos;

    public override void OnInit()
    {
        base.OnInit();
        _mainRole = this;
        _joyDir = Direction.NONE;
        isSendMove = true;

        _targetID = 0;
        _targetCollectID = 0;
        _SendMovePos = default;

        _sendMoveLast = 0;
        _sendAttackLast = 0;

        _lastSendLocation = new Point();
        _lastSendLocation.X = 0; _lastSendLocation.Y = 0;
        _lastMoveLocation = new Point { X = 0, Y = 0 };

        uID = RoleM.Instance.Uid;

        _ObjectType = ObjectType.PLAYER;

        PlayerInfo _playerInfo = GameRoomM.Instance.getPlayerInfoByID(RoleM.Instance.Uid);
        if (_playerInfo != null)
        {
            pointNow = _playerInfo.Location;

            Debug.Log("自己出生点" + pointNow.X + " pointNow.X " + pointNow.Y + " pointNow.Y");
            wrap(new Vector3((float)pointNow.X, (float)pointNow.Y, SceneManager.GetOffZ(pointNow.X, pointNow.Y)));

        }

        RefrshBuffer();
    }

    public override void RoleDeal()
    {
        base.RoleDeal();
        targetID = 0;
        targetCollectID = 0;
        _SendMovePos = default;
    }

    public override void RefrshBuffer()
    {

        _attackSpeed = ((int)RoleM.Instance.GetRoleAttrValueByAttr(pb.Attribute.ATTACK_SPEED) + 50) * 0.001f;
        moveSpeed = ((int)RoleM.Instance.GetRoleAttrValueByAttr(pb.Attribute.MOVE_SPEED) + 50) * 0.001f;
    }

    public override void AddOrRemoveBuffer(int bufferId, bool isAdd)
    {
        base.AddOrRemoveBuffer(bufferId, isAdd);

        if (bufferId == (int)CFG_BuffEffect.BufferID.CaoFeng)
        {
            if (isAdd)
            {
                if (RoleM.Instance.GetBuffOwnerUid(bufferId) != 0)
                {
                    MainRole._MainRole.targetID = RoleM.Instance.GetBuffOwnerUid(bufferId);
                    MainRole._MainRole.targetCollectID = 0;
                    MainRole._MainRole._SendMovePos = default;
                }
            }
            else
            {
                MainRole._MainRole.targetID = 0;
            }
        }
    }

    private int AutoDis = 100;
    //自动攻击检测
    private bool AutoWalkOrAttack()
    {
        if (!GameRoomM.Instance.AutoPlay && !GameRoomM.Instance.ShortAttactToPlayer)
            return false;

        if (MainRole._MainRole._SendMovePos != default)
            return false;

        if (MainRole._MainRole.targetCollectID != 0)
            return false;

        var Roles = FightController.Instance.GetCanDisRoles(AutoDis, 1);
        long TargetRole = Roles.Count > 0 ? Roles[0].role.uID : 0;


        var Monsters = MonsterManager.Instance.getCanAttackMonsters(AutoDis, 1, true);
        long TargetMonster = Monsters.Count > 0 ? Monsters[0]._monser.uID : 0;

        //Debug.Log("自动攻击检测 TargetMonster  " + TargetMonster  + "  TargetRole  " + TargetRole);
        if (TargetRole == 0 && TargetMonster == 0)
            return false;

        if (TargetMonster == 0)
            MainRole._MainRole.targetID = TargetRole;
        else if (TargetRole == 0)
            MainRole._MainRole.targetID = TargetMonster;
        else
        {
            var target = FightController.Instance.GetOneRoleInfo(TargetRole).role.pointNow;
            var monsterTarget = MonsterManager.Instance.getOneMonster(TargetMonster)._monser._gridNow;
            int _roleDis = getDistanceFromMainRole(target.X, target.Y);
            int _monsterDis = getDistanceFromMainRole(monsterTarget.X, monsterTarget.Y);

            if (_roleDis == _monsterDis)
                MainRole._MainRole.targetID = GameRoomM.Instance.AttactToPlayer ? TargetRole : TargetMonster;
            else
                MainRole._MainRole.targetID = _roleDis > _monsterDis ? TargetMonster : TargetRole;
        }
        return true;

    }

    public override void OnUpdate()
    {
        base.OnUpdate();

        if (GameRoomM.Instance._GameState != GameRoomM.GameState.Play || hp <= 0)
            return;

        _sendAttackLast += Time.deltaTime;

        _sendMoveLast += Time.deltaTime;
        if (MainRole._MainRole._joyDir != Direction.NONE)
        {
            if (_sendMoveLast > moveSpeed)
            {
                _sendMoveLast = 0;
                sendPath(_joyDir);
            }
        }
        else if (MainRole._MainRole.targetID != 0 || MainRole._MainRole.targetCollectID != 0 || MainRole._MainRole._SendMovePos != default)
        {
            if (_sendMoveLast > moveSpeed)
            {
                _sendMoveLast = 0;
                sendAutoPath();
            }
        }
        else
        {
            AutoWalkOrAttack();
        }

        //安全点不可攻击
        if (SceneManager.IsSafePos(_gridNow.X, _gridNow.Y))
            return;

        var _SkillList = SkillM.Instance.GetGameFightSKill();
        bool isNorAttack = true;
        if (checkCollection())
            isNorAttack = false;

        for (int i = 0; i < _SkillList.Count; i++)
        {
            if (checkSkillAttack(_SkillList[i]))
            {
                isNorAttack = false;
                break;
            }
        }
        if (hp > 0)
        {
            if (isNorAttack)
                checkNormalAttack();
        }
    }

    public override void moveTo(Point Location, Direction _Direction)
    {

        base.moveTo(Location, _Direction);
        _lastMoveLocation = Location;

        if (MainRole._MainRole._SendMovePos != default && MainRole._MainRole._SendMovePos.X == Location.X && MainRole._MainRole._SendMovePos.Y == Location.Y)
            MainRole._MainRole._SendMovePos = default;
    }

    private void OnDisable() { }

    //根据行动目标寻路移动
    public bool sendAutoPath()
    {
        if (_sendAttackLast < _attackSpeed)
            return false;

        if (hp <= 0)
            return false;

        Point toPoint = pointNow;
        if (MainRole._MainRole._SendMovePos == default)
        {
            long _toID = MainRole._MainRole.targetCollectID == 0 ? MainRole._MainRole.targetID : MainRole._MainRole.targetCollectID;

            if (_toID == 0)
            {
                return false;
            }

            //攻击距离
            float attackDis = (float)CFG_Skill.GetById(GameConst.GetNormallSkillByJob(job)).MaxDistance;//float attackDis = GameRoomM.Instance.GetCanNoAttackGame() ? 1 : (float)CFG_Skill.GetById(GameConst.GetNormallSkillByJob(job)).MaxDistance;

            var target = FightController.Instance.GetOneRoleInfo(_toID);
            if (target != null)
            {
                toPoint = target.role.pointNow;
            }
            else
            {
                var monsterTarget = MonsterManager.Instance.getOneMonster(_toID);
                if (monsterTarget != null)
                {
                    toPoint = monsterTarget._monser._gridNow;

                    //如果是采集怪物则攻击距离为1
                    if (monsterTarget._monser.maxLifeTime > 0)
                        attackDis = 1;
                }
            }

            if (toPoint == pointNow)
            {
                return false;
            }

            //检测是否在攻击距离内（在攻击范围内不移动）
            int Dis = getDistanceFromMainRole(toPoint.X, toPoint.Y);

            //Debug.Log($"根据攻击目标寻路移动  Dis {Dis} AutoDis {AutoDis} ");
            //非嘲讽状态的寻路 超过指定距离就不再检测
            if (!RoleM.Instance.CheakBuffId((int)CFG_BuffEffect.BufferID.CaoFeng) && Dis > AutoDis)
                return false;

            if (attackDis >= Dis)
            {
                return false;
            }

        }
        else
        {
            toPoint = MainRole._MainRole._SendMovePos;

            if (toPoint == pointNow)
            {
                MainRole._MainRole._SendMovePos = default;
                return false;
            }
        }

        Point posNow = pointNow;
        double x = posNow.X;
        double y = posNow.Y;

        var offx = toPoint.X - pointNow.X;
        var offy = toPoint.Y - pointNow.Y;

        if (Math.Abs(offx) > 0)
        {
            if (offx > 0)
                x += MoveOff;
            else
                x -= MoveOff;
        }
        if (Math.Abs(offy) > 0)
        {
            if (offy > 0)
                y += MoveOff;
            else
                y -= MoveOff;
        }

        if (x == _lastMoveLocation.X && y == _lastMoveLocation.Y)
            return false;


        if (SceneManager.GetOffZ(x, y) != SceneManager.Block)
        {
            _lastSendLocation.X = x;
            _lastSendLocation.Y = y;
            TcpC.C2S_Move(new MoveRequest
            {
                Location = _lastSendLocation,
            });

        }
        return true;
    }

    private float MoveOff = 1f;
    public void sendPath(Direction dir)
    {
        //if (_Actionset > 0)
        //{
        //    return;
        //}

        Point posNow = pointNow;

        double x = posNow.X;
        double y = posNow.Y;
        switch (dir)
        {
            case Direction.NONE:

                break;
            case Direction.UP:
                y += MoveOff;
                break;
            case Direction.DOWN:
                y -= MoveOff;
                break;
            case Direction.LEFT:
                x -= MoveOff;
                break;
            case Direction.RIGHT:
                x += MoveOff;
                break;
            case Direction.UP_LEFT:
                x -= MoveOff;
                y += MoveOff;
                break;
            case Direction.UP_RIGHT:
                x += MoveOff;
                y += MoveOff;
                break;
            case Direction.DOWN_LEFT:
                x -= MoveOff;
                y -= MoveOff;
                break;
            case Direction.DOWN_RIGHT:
                x += MoveOff;
                y -= MoveOff;
                break;
        }

        if (x == _lastMoveLocation.X && y == _lastMoveLocation.Y)
            return;

        if (hp <= 0)
            return;

        if (SceneManager.GetOffZ(x, y) != SceneManager.Block)
        {
            //_sendAttackLast = 0;
            //_Actionset = 0;
            //Debug.Log($"发送移动坐标 {x} {y} ");

            _lastSendLocation.X = x;
            _lastSendLocation.Y = y;

            TcpC.C2S_Move(new MoveRequest
            {
                Location = _lastSendLocation,
            });
        }

    }

    public int getDistanceFromMainRole(double x, double y)
    {
        //Vector2.Distance(new Vector2((float)pointNow.X, (float)pointNow.Y), new Vector2(x, y));
        int cutX = (int)Math.Abs(x - pointNow.X);
        int cutY = (int)Math.Abs(y - pointNow.Y);
        return Math.Max(cutX, cutY);
    }

    //技能攻击检测条件判定
    private bool AttackCheak(int skillId = 0)
    {
        long Pid = skillId == 0 ? 0 : SkillM.Instance.GetUiDBySkillId(skillId);

        //if (_Actionset > 0 && Pid == 0)
        //{
        //    return false;
        //}
        //普攻检测自己
        if (Pid == 0 && (hp <= 0 || IsMoving()) || MainRole._MainRole.targetCollectID != 0 || MainRole._MainRole._SendMovePos != default)
            return false;
        else if (Pid != 0)
        {
            //自己伙伴的技能 检测伙伴的血量
            if (FightController.Instance.GetOneRoleInfo(Pid) == null || FightController.Instance.GetOneRoleInfo(Pid).role.hp <= 0)
                return false;
        }

        if (GameRoomM.Instance.GetCanNoAttackGame())
            return false;

        if (GameRoomM.Instance._GameState != GameRoomM.GameState.Play)
            return false;

        if (Pid == 0 && RoleM.Instance.CheakBuffId((int)CFG_BuffEffect.BufferID.XuanYun, (int)CFG_BuffEffect.BufferID.BingDong))
            return false;

        if (Pid == 0 && _sendAttackLast < _attackSpeed)
            return false;

        return true;
    }

    //技能攻击检测
    private bool checkSkillAttack(int SkillID)
    {
        if (!SkillM.Instance.GetGameSkillOpen(SkillID))
            return false;

        if (SkillM.Instance.GetGameSkillDealy(SkillID) >= GameUtil.Instance.GetCurrentTimeSeconds())
            return false;

        if (!AttackCheak(SkillID))
            return false;

        var CFGSKill = CFG_Skill.GetById(SkillID);
        if (CFGSKill.PreMark == 1)
            return false;

        //if (CFGSKill.Skilltype == (int)CFG_Skill.SkillType.Super || CFGSKill.Skilltype == (int)CFG_Skill.SkillType.SuperCtr)
        //{
        //    SkillM.Instance.SetGameSkillOpen(SkillID, false, 0);
        //}

        bool IsNoTarget = CFGSKill.Needaim == 0;
        Direction dir = _Dir;

        long attacktargetID;
        Point targetPos = null;

        if (!IsNoTarget)
        {
            attacktargetID = GetAttackUid(SkillID, ref targetPos);
            if (attacktargetID == 0)
                return false;

            MainRole._MainRole.targetID = attacktargetID;
            dir = getDir(targetPos);
        }
        else
        {
            attacktargetID = 0;
        }

        //Log.PrintYellow("SkillAttack targetID:" + targetID + " dir:" + dir + "Time" + GameUtil.Instance.GetCurrentMSTimeSeconds()
        //    + "role pos is:" + pointNow.X + " " + pointNow.Y + "SkillID  " + SkillID);

        TcpC.C2S_Attack(new AttackRequest
        {
            SkillId = SkillID,
            TargetUids = new long[] { attacktargetID },
            Direction = dir,
            Location = IsNoTarget ? pointNow : targetPos,
        });

        if (SkillM.Instance.GetUiDBySkillId(SkillID) != 0)
        {
            return false;
        }

        _sendAttackLast = 0;
        //_sendMoveLast = 0;

        return true;
    }

    /// <summary>
    /// 根据技能ID 获取攻击目标UID
    /// </summary>
    /// <param name="SkillId"></param>
    /// <param name="point"></param>
    /// <param name="isAttack"></param>
    /// <returns></returns>
    private long GetAttackUid(int SkillId, ref Point point)
    {
        int attackDis = (int)CFG_Skill.GetById(SkillId).MaxDistance;
        long mId = MonsterManager.Instance.getCanAttackMonster(SkillId, MainRole._MainRole.targetID, attackDis, true);
        if (mId != 0 && mId == MainRole._MainRole.targetID)
        {
            point = MonsterManager.Instance.getMonsterPos(mId);
            return mId;
        }

        long uId = FightController.Instance.getCanAttackMonster(SkillId, MainRole._MainRole.targetID, attackDis);
        if (uId != 0 && (uId == MainRole._MainRole.targetID || GameRoomM.Instance.AttactToPlayer))
        {
            point = FightController.Instance.getRolePos(uId);
            return uId;

        }
        if (mId > 0)
        {
            point = MonsterManager.Instance.getMonsterPos(mId);
            return mId;
        }
        else if (uId > 0)
        {
            point = FightController.Instance.getRolePos(uId);
            return uId;
        }
        return 0;
    }

    //普攻检测
    private void checkNormalAttack()
    {
        if (!AttackCheak())
            return;

        Point monsterPos = default;

        long targetID = GetAttackUid(GameConst.GetNormallSkillByJob(job), ref monsterPos);
        if (targetID == 0)
            return;
        MainRole._MainRole.targetID = targetID;
        _sendAttackLast = 0;
        //_sendMoveLast = 0;


        //Log.Print("j普攻检测:------------+" + _Actionset + " nowTime" + GameUtil.Instance.GetCurrentMSTimeSeconds());

        Direction dir = getDir(monsterPos);

        TcpC.C2S_Attack(new AttackRequest
        {
            SkillId = GameConst.GetNormallSkillByJob(job),
            TargetUids = new long[] { targetID },
            Direction = dir
        });
    }

    //可收集 怪物检测
    private long GetCollectionUid(int SkillId)
    {
        if (MainRole._MainRole.targetCollectID != 0)
            return MainRole._MainRole.targetCollectID;

        long mId = MonsterManager.Instance.getCanAttackMonster(0, MainRole._MainRole.targetCollectID, 1, false);
        if (mId != 0)
        {
            MainRole._MainRole.targetCollectID = mId;
            return mId;
        }

        return 0;
    }

    //采集检测
    private bool checkCollection()
    {
        if (_sendAttackLast < _attackSpeed)
            return false;

        long targetID = MainRole._MainRole.targetCollectID;
        //采集副本时候可以搜索可采集的怪物
        if (GameRoomM.Instance.GetCanNoAttackGame())
        {
            targetID = GetCollectionUid(0);
        }
        else
        {
            //非采集怪时不用检测采集
            var _Monster = MonsterManager.Instance.getOneMonster(targetID);
            if (_Monster == null || _Monster._monser.maxLifeTime <= 0)
                return false;
        }

        if (targetID == 0)
            return false;

        var Monster = MonsterManager.Instance.getOneMonster(targetID);
        if (Monster == null || getDistanceFromMainRole(Monster._monser._gridNow.X, Monster._monser._gridNow.Y) > 1)
            return false;


        if (Monster._monser.curLifeTime == 0)
        {
            TcpC.C2S_Mining(targetID);
        }

        Monster._monser.curLifeTime += _attackSpeed;
        _sendAttackLast = 0;

        if (Monster._monser.curLifeTime >= Monster._monser.maxLifeTime)
        {
            Monster._monser.curLifeTime = Monster._monser.maxLifeTime;
        }
        return true;

    }


    public override int getHurtDistance(int skillID, int dis = -1, bool isAttactRole = false)
    {
        return -1;
    }
    //private int EffIdx = 0;

    public void SetDropEff(Point effPos, List<ItemInfo> ItemInfos)
    {
        Vector3 toPos = new Vector3((float)effPos.X, (float)effPos.Y, SceneManager.EffOffZ);

        //foreach (var item in ItemInfos)
        //{
        //    if (CFG_Item.GetById(item.CfgId).Quality > 2)
        //    {

        //        string poolKey = "drop_" + CFG_Item.GetById(item.CfgId).Quality;
        //        if (!PoolHelper.Instance.CheckPool(poolKey))
        //        {
        //            var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
        //            await baseObj.WaitForAsyncLoading();
        //            PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
        //        }

        //        var Eff = PoolHelper.Instance.GetPoolObj(poolKey);
        //        Eff.transform.SetParent(FightController.Instance.Root_Eff3D);
        //        Eff.transform.position = toPos;
        //        Eff.SetActive(true);
        //        string effJacKey = "DropEff" + EffIdx;

        //        Eff.transform.DOMove(this.transform.position, 0.2f).SetDelay(1f).OnComplete(() =>
        //        {
        //            PoolHelper.Instance.RecycObj(Eff);
        //        });

        //        EffIdx++;
        //    }

        //}
    }

    public void OnDestroy()
    {
        //_mainRole = null;
    }
}

