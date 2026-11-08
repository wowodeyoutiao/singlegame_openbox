using Cysharp.Threading.Tasks;
using DG.Tweening;
using HotUpdateScripts;
using HotUpdateScripts.Game;
using HotUpdateScripts.ProtoM;
using JEngine.Core;
using JEngine.Helper;
using pb;
using Spine.Unity;
using System;
using System.Collections.Generic;
using UnityEngine;

public class Role : MonoBehaviour
{

    public SkeletonAnimation anim;
    public GameConst.ModelAnim _animState;

    //public NavMeshAgent agent;
    public Point _gridNow;//当前格子

    private int _maxHp;//最大血量
    public int maxHp
    {
        get => _maxHp;

        set
        {
            if (_maxHp == value)
                return;
            _maxHp = value;
        }
    }


    private int _hpNow;//当前血量
    public int hp
    {
        get => _hpNow;

        set
        {
            if (_hpNow == value)
                return;
            _hpNow = value;

        }
    }

    private int _maxMp;//最大蓝
    public int maxMp
    {
        get => _maxMp;

        set
        {
            if (_maxMp == value)
                return;
            _maxMp = value;

        }
    }

    private int _mpNow;//当前蓝
    public int mpNow
    {
        get => _mpNow;

        set
        {
            if (_mpNow == value)
                return;
            _mpNow = value;

        }
    }


    private SexType _sex;//性别
    public SexType sex
    {
        get => _sex;

        set
        {
            if (_sex == value)
                return;
            _sex = value;

        }
    }

    private JobType _job;//职业
    public JobType job
    {
        get => _job;

        set
        {
            if (_job == value)
                return;
            _job = value;

        }
    }

    private int _level;//等级
    public int level
    {
        get => _level;

        set
        {
            if (_level == value)
                return;
            _level = value;
        }
    }

    private string _roleName;//玩家名字
    public string roleName
    {
        get => _roleName;

        set
        {
            _roleName = value;

        }
    }

    private long _uID;//玩家ID
    public long uID
    {
        get => _uID;

        set
        {
            if (_uID == value)
                return;
            _uID = value;

        }
    }

    private string _headUrl;//头像地址
    public string HeadUrl
    {
        get => _headUrl;

        set
        {
            _headUrl = value;
        }
    }

    private Point _pointNow;
    public Point pointNow
    {
        get => _pointNow;

        set
        {
            if (_pointNow == value)
                return;
            _pointNow = value;

        }
    }

    private bool _isSendMove;
    public bool isSendMove
    {
        get => _isSendMove;

        set
        {
            if (_isSendMove == value)
                return;
            _isSendMove = value;

        }
    }

    private float _moveSpeed;
    public float moveSpeed
    {
        get => _moveSpeed;

        set
        {
            if (_moveSpeed == value)
                return;
            _moveSpeed = value;

        }
    }

    private Direction _dir;
    public Direction _Dir
    {
        get => _dir;

        set
        {
            if (_dir == value)
                return;
            _dir = value;

        }
    }

    //类型
    private ObjectType _objectType;
    public ObjectType _ObjectType
    {
        get => _objectType;
        set
        {
            _objectType = value;
        }
    }

    /// <summary>
    /// 归属ID 用于伙伴的主角判定
    /// </summary>
    public long OwnerUid
    {
        get; set;
    }

    /// <summary>
    /// 用于伙伴的CFG_ID
    /// </summary>
    public int CfgID { get; set; }

    /// <summary>
    /// 行会ID
    /// </summary>
    public long GuildId { get; set; }

    /// <summary>
    /// 行会名称
    /// </summary>
    public string GuildName { get; set; }

    /// <summary>
    /// 是否敌对角色
    /// </summary>
    public bool isEnemy = true;

    //称号ID
    public int TitleId { get; set; } = 0;


    public Transform root_hub;
    private Dictionary<string, GameObject> ObjDic = new Dictionary<string, GameObject>();
    public List<int> _BufferList = new List<int>();
    public void AddJaction(string key, GameObject jac)
    {
        if (!ObjDic.ContainsKey(key))
        {
            ObjDic.Add(key, jac);
        }
        else
        {
            ClearOneObj(key);
            ObjDic[key] = jac;
        }
    }

    public void CleacObjDic()
    {
        if (PoolHelper.Instance == null)
            return;

        foreach (var item in ObjDic.Values)
        {
            PoolHelper.Instance.RecycObj(item);
        }
        ObjDic.Clear();
    }

    public void ClearOneObj(string key)
    {
        if (ObjDic.ContainsKey(key))
        {
            PoolHelper.Instance.RecycObj(ObjDic[key]);
            ObjDic.Remove(key);
        }
    }

    public virtual async void AddOrRemoveBuffer(int bufferId, bool isAdd = true)
    {
        if (isAdd)
        {
            if (!_BufferList.Contains(bufferId))
                _BufferList.Add(bufferId);
        }
        else
        {
            if (_BufferList.Contains(bufferId))
                _BufferList.Remove(bufferId);
        }
        if (!_isShowEff) return;

        if (GameRoomM.Instance.GameAttack == GameRoomM.AttackClass.Team
           && OwnerUid != MainRole._MainRole.uID)
            return;
        //Log.Print("bufferId   " + bufferId);
        string effKey = CFG_BuffEffect.GetById(bufferId).BuffEffect;
        if (!string.IsNullOrEmpty(effKey))
        {
            if (isAdd)
            {
                string poolKey = effKey;
                if (!PoolHelper.Instance.CheckPool(poolKey))
                {
                    var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
                    await UniTask.WaitUntil(() => baseObj.Loaded);
                    PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
                }
                await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
                GameObject effObj = PoolHelper.Instance.GetPoolObj(effKey);
                effObj.transform.SetParent(this.transform);
                effObj.transform.localPosition = Vector3.back * 0.1f;
                effObj.gameObject.name = effKey;
                effObj.SetActive(true);
                AddJaction(effKey, effObj);
            }
            else
            {
                ClearOneObj(effKey);
            }
        }

    }

    public void OnInitBuffer(List<BuffInfo> _buffInfos)
    {
        _BufferList.Clear();
        if (_buffInfos != null)
        {
            int Idx = 1;
            foreach (BuffInfo item in _buffInfos)
            {
                string jkey = "InitBuffer_" + item.BuffId;
                AddOrRemoveBuffer(item.BuffId, true);
                _BufferList.Add(item.BuffId);
                Idx++;
            }
        }

        _gridNow = new Point();

        PlayerInfo _playerInfo = GameRoomM.Instance.getPlayerInfoByID(uID);
        Direction direction = default;
        if (_playerInfo != null)
        {
            pointNow = _playerInfo.Location;
            direction = _playerInfo.Direction;
        }
        else
        {
            ObjectInfo objectInfo = GameRoomM.Instance.GetOneObjInfoByID(uID);
            if (objectInfo != null)
            {
                pointNow = objectInfo.Location;
                direction = objectInfo.Direction;
            }
        }

        wrap(new Vector3((float)pointNow.X, (float)pointNow.Y, SceneManager.GetOffZ(pointNow.X, pointNow.Y)));
        setRotation(direction);

        if (hp <= 0)
        {
            RoleDeal();
        }

        isEnemy = !(OwnerUid == MainRole._MainRole.uID || (0 != MainRole._MainRole.GuildId && GuildId == MainRole._MainRole.GuildId));
    }

    public virtual void RefrshBuffer()
    {
    }

    public virtual void OnInit()
    {
        if (this.transform.Find("spine"))
            anim = this.transform.Find("spine").GetComponent<SkeletonAnimation>();
        if (this.transform.Find("root_hub"))
            root_hub = this.transform.Find("root_hub");
        _gridNow = new Point();

        moveSpeed = 0.2f;
        _Dir = Direction.LEFT;
        _ObjectType = ObjectType.PLAYER;

        OwnerUid = 0;

        _animState = GameConst.ModelAnim.No;
        changeState(GameConst.ModelAnim.stand);
    }

    public void changeState(GameConst.ModelAnim actName)
    {
        if (actName != GameConst.ModelAnim.dead && actName != GameConst.ModelAnim.attack && !_animState.ToString().Contains("skill") && _animState.CompareTo(actName) == 0)
            return;
        _animState = actName;

        if (anim == null || root_hub == null || (_ObjectType != ObjectType.PLAYER && _ObjectType != ObjectType.PARTNER))
            return;

        //Debug.Log("play  :" + actName.ToString() + " Uid" + uID);
        if (_animState == GameConst.ModelAnim.attack || _animState.ToString().Contains("skill"))
            SpineUtil.Instance.PlayOnce(anim, _animState.ToString(), () =>
            {
                _animState = GameConst.ModelAnim.stand;
                SpineUtil.Instance.PlayLoop(anim, _animState.ToString());
            }, 3f);
        else if (actName == GameConst.ModelAnim.dead)
            SpineUtil.Instance.PlayOnce(anim, _animState.ToString(), null);
        else
            SpineUtil.Instance.PlayLoop(anim, _animState.ToString());
    }

    private bool _isShowEff = true;

    private HashSet<string> EffHasSet = new HashSet<string>();
    public async UniTaskVoid SetHurtChange(int changeNum, int MaxHp, long AttackerId, int skillId)
    {
        hp += changeNum;
        maxHp = MaxHp;

        if (GameRoomM.Instance.GameAttack == GameRoomM.AttackClass.Team && OwnerUid != MainRole._MainRole.uID && !GameRoomM.Instance.CheakIsMainRole(AttackerId))
            return;

        if (changeNum < 0)
        {
            string poolKey = "Hurt_1";

            if (skillId > 0 && CFG_Skill.GetById(skillId).SkillProfession % 10 != 0)
            {
                poolKey = "Hurt_" + (CFG_Skill.GetById(skillId).SkillProfession % 10);
            }

            if (EffHasSet.Contains(poolKey))
                return;
            //else if (skillId > 0 && MonsterManager.Instance.getOneMonster(AttackerId) != null)
            //{
            //    poolKey = CFG_Monster.GetById(MonsterManager.Instance.getOneMonster(AttackerId)._monser.cfgId).AttackDistance > 1 ? "Hurt_2" : "Hurt_1";
            //}
            AudioManager.Instance.PlaySound((AudioKey.Sound)Enum.Parse(typeof(AudioKey.Sound), poolKey));

            if (!PoolHelper.Instance.CheckPool(poolKey))
            {
                var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
                await UniTask.WaitUntil(() => baseObj.Loaded);
                PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
            }
            await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
            var Eff = PoolHelper.Instance.GetPoolObj(poolKey);
            Eff.transform.SetParent(FightController.Instance.Root_Eff3D);
            Eff.transform.position = new Vector3(transform.position.x, transform.position.y, SceneManager.EffOffZ); ;
            Eff.SetActive(true);
            EffHasSet.Add(poolKey);
            CoroutineCall(0.8f, () =>
            {
                EffHasSet.Remove(poolKey);
                if (Eff != null)
                    PoolHelper.Instance.RecycObj(Eff);
            });
        }
    }

    //private void CoroutineCall(float time, Action callBack)
    //{
    //    UniCoroutineCall(time, callBack).Forget();
    //}

    //private async UniTaskVoid UniCoroutineCall(float time, Action callBack)
    //{
    //    await UniTask.Delay((int)(time * 1000), cancellationToken: this.GetCancellationTokenOnDestroy());
    //    callBack?.Invoke();
    //}

    private void ClearCoroutine()
    {
        if (CoroutineMgr.Instance == null)
            return;

        foreach (var item in coroutinesDic)
        {
            CoroutineMgr.Instance.StopCoroutine(item.Value);
        }
        coroutinesDic.Clear();
    }
    private int coroutineId = 0;
    private Dictionary<int, Coroutine> coroutinesDic = new Dictionary<int, Coroutine>();
    private Coroutine CoroutineCall(float time, Action callBack)
    {
        var coroutineItem = CoroutineMgr.Instance.StartCoroutine(IECoroutineCall(time, callBack, coroutineId));
        coroutinesDic.Add(coroutineId, coroutineItem);
        coroutineId++;
        return coroutineItem;
    }


    private System.Collections.IEnumerator IECoroutineCall(float time, Action callBack, int cID)
    {
        yield return new WaitForSeconds(time);
        callBack?.Invoke();
        coroutinesDic.Remove(cID);
    }

    public async UniTaskVoid SetUpLv()
    {
        if (GameRoomM.Instance.GameAttack == GameRoomM.AttackClass.Team
         && OwnerUid != MainRole._MainRole.uID)
            return;

        string poolKey = "LvUp";
        if (!PoolHelper.Instance.CheckPool(poolKey))
        {
            var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
            await UniTask.WaitUntil(() => baseObj.Loaded);
            PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
        }
        await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
        var Eff = PoolHelper.Instance.GetPoolObj(poolKey);
        Eff.transform.SetParent(FightController.Instance.Root_Eff3D);
        Eff.transform.position = new Vector3(transform.position.x, transform.position.y, SceneManager.EffOffZ); ;
        Eff.SetActive(true);

        CoroutineCall(1.36f, () =>
        {
            if (Eff != null)
                PoolHelper.Instance.RecycObj(Eff);
        });
    }

    private void Update()
    {
        OnUpdate();
    }

    public virtual void OnUpdate()
    {
        if (_animState == GameConst.ModelAnim.walk || _animState == GameConst.ModelAnim.run)
        {
            if (!IsMoving())
            {
                changeState(GameConst.ModelAnim.stand);
            }
        }
    }

    //private float requiredSpeed = 0;
    private Vector3 _runToPos;
    //private float _runToTime = 0.2f; //移动时间

    public virtual void moveTo(Point Location, Direction _Direction)
    {
        if (_animState == GameConst.ModelAnim.dead)
            return;

        isSendMove = true;
        float OffZ = SceneManager.GetOffZ(Location.X, Location.Y);
        if (OffZ != SceneManager.Block)
        {
            //Debug.Log("jerrycrole " + uID + "server set des to " + Location.X + "  " + Location.Y
            //    +"  time:" + DateTime.Now.TimeOfDay);
            _runToPos = new Vector3((float)Location.X, (float)Location.Y, OffZ);
            //_runToTime = moveSpeed + 0.2f;
            //requiredSpeed = Vector3.Distance(transform.position, _runToPos) / _runToTime;

            changeState(GameRoomM.Instance.CheakHaveBuffer(uID, (int)CFG_BuffEffect.BufferID.AttackAndMoveAdd) ? GameConst.ModelAnim.run : GameConst.ModelAnim.walk);
            isEndMove = true;

            this.transform.DOKill();
            this.transform.DOMove(_runToPos, (moveSpeed + 0.2f)).SetEase(Ease.Linear).OnComplete(() =>
            {
                isEndMove = false;
            });
            setRotation(_Direction);

            _gridNow.X = Math.Floor(Location.X);
            _gridNow.Y = Math.Floor(Location.Y);

            pointNow.X = Math.Floor(Location.X);
            pointNow.Y = Math.Floor(Location.Y);
        }

    }

    public virtual void wrap(Vector3 target)
    {
        this.transform.DOKill();
        this.transform.position = target;
        _gridNow.X = target.x;
        _gridNow.Y = target.y;
    }

    public bool isEndMove = false;
    public bool IsMoving()
    {
        return isEndMove;
    }
    public virtual void destoryRole()
    {
        CleacObjDic();
        ClearCoroutine();
        string pooKey = this.gameObject.name.Split('/')[0];
        PoolHelper.Instance.RecycObj(pooKey, this.gameObject);
    }

    public virtual void RoleDeal()
    {
        //Debug.Log("RoleDeal");
        hp = 0;
        isEndMove = false;
        this.transform.DOKill();
        changeState(GameConst.ModelAnim.dead);
        CleacObjDic();
    }

    void OnDestroy()
    {
        CleacObjDic();
        ClearCoroutine();
    }

    public async UniTaskVoid RoleRevive(int HP, Point Location)
    {
        hp = HP;
        isSendMove = true;
        changeState(GameConst.ModelAnim.stand);
        pointNow = Location;
        wrap(new Vector3((float)Location.X, (float)Location.Y, SceneManager.GetOffZ(Location.X, Location.Y)));

        if (_BufferList != null)
        {
            foreach (int BuffId in _BufferList)
            {
                string jkey = "InitBuffer_" + BuffId;
                AddOrRemoveBuffer(BuffId, true);
            }
        }

        if (GameRoomM.Instance.GameAttack == GameRoomM.AttackClass.Team
        && OwnerUid != MainRole._MainRole.uID)
            return;

        string poolKey = "eff_HuiGuangFanZao";
        if (!PoolHelper.Instance.CheckPool(poolKey))
        {
            var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
            await UniTask.WaitUntil(() => baseObj.Loaded);
            PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
        }
        await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
        var Eff = PoolHelper.Instance.GetPoolObj(poolKey);
        Eff.transform.SetParent(FightController.Instance.Root_Eff3D);
        Eff.transform.position = new Vector3(transform.position.x, transform.position.y, SceneManager.EffOffZ);
        Eff.SetActive(true);
        CoroutineCall(1, () =>
        {
            if (Eff != null)
                PoolHelper.Instance.RecycObj(Eff);
        });
    }

    public void setRotation(Direction dir)
    {
        //Log.Print("jerryc Move setRotation:" + dir);
        this.transform.eulerAngles = SceneManager.GetRotationPos(dir, ref _dir);
    }

    //传送
    public async UniTaskVoid SetTeleport(Point point, Direction dir, TeleportType Type)
    {
        wrap(new Vector3((float)point.X, (float)point.Y, SceneManager.GetOffZ(point.X, point.Y)));
        setRotation(dir);
        pointNow = point;
        if (Type == TeleportType.BLINK)
        {
            string poolKey = "eff_ChuanSong";
            if (!PoolHelper.Instance.CheckPool(poolKey))
            {
                var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
                await UniTask.WaitUntil(() => baseObj.Loaded);
                PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
            }
            await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));

            var Eff = PoolHelper.Instance.GetPoolObj(poolKey);
            Eff.transform.SetParent(FightController.Instance.Root_Eff3D);
            Eff.transform.position = new Vector3(transform.position.x, transform.position.y, SceneManager.EffOffZ);
            Eff.transform.eulerAngles = SceneManager.GetRotationPos(dir, ref _dir);
            Eff.SetActive(true);
            CoroutineCall(1, () =>
            {
                if (Eff != null)
                    PoolHelper.Instance.RecycObj(Eff);
            });
        }
    }

    public async UniTaskVoid SetAttack(int SkillId, long[] TargetUids, Point Location, Direction attackDirection)
    {
        var CFGSKIll = CFG_Skill.GetById(SkillId);
        GameConst.ModelAnim animName = GameConst.GetAnim((int)CFGSKIll.Needskill);

        changeState(animName);
        if (!_isShowEff) return;

        if (CFG_Actionset.ChealkId(CFGSKIll.SkillID) && !string.IsNullOrEmpty(CFG_Actionset.GetById(CFGSKIll.SkillID).Music))
        {
            AudioManager.Instance.PlaySound(CFG_Actionset.GetById(CFGSKIll.SkillID).Music);
        }

        if (!CFG_SkillEffect.Get().DicDatas.ContainsKey(CFGSKIll.SkillID))
            return;

        if (GameRoomM.Instance.GameAttack == GameRoomM.AttackClass.Team
            && OwnerUid != MainRole._MainRole.uID && (TargetUids == null || !GameRoomM.Instance.CheakIsMainRole(TargetUids[0])))
            return;

        var skillEff = CFG_SkillEffect.GetById(CFGSKIll.SkillID);
        Vector3 tpos;
        if (Location == null)
        {
            if (TargetUids == null || TargetUids.Length <= 0)
            {
                return;
            }

            if (FightController.Instance.GetOneRoleInfo(TargetUids[0]) != null)
            {
                tpos = FightController.Instance.GetOneRoleInfo(TargetUids[0]).obj.transform.position;
            }
            else if (MonsterManager.Instance.getOneMonster(TargetUids[0]) != null)
            {
                tpos = MonsterManager.Instance.getOneMonster(TargetUids[0])._monsterObj.transform.position;
            }
            else
            {
                return;
            }
        }
        else
        {
            tpos = new Vector3((float)Location.X, (float)Location.Y, SceneManager.GetOffZ(Location.X, Location.Y));
        }
        tpos.z = SceneManager.EffOffZ;

        if (!string.IsNullOrEmpty(skillEff.BallisticEffect))
        {
            string poolKey = skillEff.BallisticEffect;
            if (!PoolHelper.Instance.CheckPool(poolKey))
            {
                var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
                await UniTask.WaitUntil(() => baseObj.Loaded);
                PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
            }
            await UniTask.WaitUntil(() => PoolHelper.Instance != null && PoolHelper.Instance.CheckPool(poolKey));
            var Eff = PoolHelper.Instance.GetPoolObj(poolKey);
            Eff.transform.SetParent(FightController.Instance.Root_Eff3D);
            Eff.transform.position = new Vector3(transform.position.x, transform.position.y, SceneManager.EffOffZ);

            Vector2 direction = (transform.position - tpos);
            direction.Normalize();
            // 设置物体的朝向
            Eff.transform.rotation = Quaternion.AngleAxis(Mathf.Atan2(direction.y, direction.x) * Mathf.Rad2Deg, Vector3.forward);

            CoroutineCall((float)skillEff.BallisticDelay / 1000f, () =>
            {
                if (Eff == null)
                    return;
                Eff.SetActive(true);
                Eff.transform.DOMove(tpos, (float)skillEff.BallisticRemainTime / 1000).OnComplete(() =>
                {
                    Eff.transform.DOKill();
                    PoolHelper.Instance.RecycObj(Eff);
                });
            });

        }

        if (!string.IsNullOrEmpty(skillEff.StartEffect))
        {
            string poolKey = skillEff.StartEffect;
            if (!PoolHelper.Instance.CheckPool(poolKey))
            {
                var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
                await UniTask.WaitUntil(() => baseObj.Loaded);
                PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
            }
            await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
            var Eff = PoolHelper.Instance.GetPoolObj(poolKey);
            Eff.transform.SetParent(FightController.Instance.Root_Eff3D);
            CoroutineCall((float)skillEff.StartDelay / 1000f, () =>
            {
                if (Eff == null)
                    return;
                Eff.transform.position = new Vector3(transform.position.x, transform.position.y, SceneManager.EffOffZ);

                if (skillEff.isAmendStart != 0)
                    Eff.transform.eulerAngles = SceneManager.GetSeverRotationPos(attackDirection);
                else
                    Eff.transform.eulerAngles = SceneManager.GetRotationPos(_Dir, ref _dir);

                Eff.SetActive(true);
            });

            CoroutineCall((float)skillEff.StartRemainTime / 1000f, () =>
            {
                if (Eff == null)
                    return;
                PoolHelper.Instance.RecycObj(Eff);
            });
        }

        if (!string.IsNullOrEmpty(skillEff.BodyEffect))
        {
            string poolKey = skillEff.BodyEffect;
            if (!PoolHelper.Instance.CheckPool(poolKey))
            {
                var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
                await UniTask.WaitUntil(() => baseObj.Loaded);
                PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
            }
            await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));

            var Eff = PoolHelper.Instance.GetPoolObj(poolKey);
            Eff.transform.SetParent(FightController.Instance.Root_Eff3D);
            Eff.transform.position = tpos;

            CoroutineCall((float)skillEff.BodyDelay / 1000f, () =>
            {
                if (Eff == null)
                    return;
                Eff.SetActive(true);
                CoroutineCall((float)skillEff.BodyRemainTime / 1000f, () =>
                {
                    if (Eff == null)
                        return;
                    PoolHelper.Instance.RecycObj(Eff);
                });
            });

        }
    }

    public Direction getDir(Point targetPoint)
    {
        Direction dir = Direction.UP;

        if (pointNow.X == targetPoint.X)
        {
            if (pointNow.Y > targetPoint.Y)
                dir = Direction.DOWN;
            else
                dir = Direction.UP;

        }
        else if (pointNow.X > targetPoint.X)
        {

            if (pointNow.Y == targetPoint.Y)
                dir = Direction.LEFT;
            else if (pointNow.Y < targetPoint.Y)
                dir = Direction.DOWN_LEFT;
            else
                dir = Direction.UP_LEFT;

        }
        else if (pointNow.X < targetPoint.X)
        {

            if (pointNow.Y == targetPoint.Y)
                dir = Direction.RIGHT;
            else if (pointNow.Y > targetPoint.Y)
                dir = Direction.DOWN_RIGHT;
            else
                dir = Direction.UP_RIGHT;

        }

        return dir;
    }

    public virtual int getHurtDistance(int skillID, int dis = -1, bool isAttactRole = false)
    {
        if (OwnerUid == MainRole._MainRole.uID)
            return -1;

        if (hp <= 0)
            return -1;

        if (_animState == GameConst.ModelAnim.dead)
            return -1;

        if (_ObjectType != ObjectType.PLAYER && _ObjectType != ObjectType.PARTNER)
        {
            return -1;
        }
        //安全点不可攻击
        if (SceneManager.IsSafePos(_gridNow.X, _gridNow.Y))
            return -1;

        if (_ObjectType != ObjectType.PLAYER && isAttactRole)
        {
            return -1;
        }

        //工会成员互相不可攻击
        if (GuildId != 0 && GuildId == MainRole._MainRole.GuildId)
            return -1;

        //int distance = MainRole._MainRole.getDistanceFromMainRole(_gridNow.X, _gridNow.Y);
        int distance = FightController.Instance.getDistance(_gridNow.X, _gridNow.Y, SkillM.Instance.GetUiDBySkillId(skillID));

        if (dis > 0)
        {
            if (dis >= distance)
            {
                return (int)distance;
            }
        }
        else
        {
            var cfgInfo = CFG_Skill.GetById(skillID);
            decimal needDistance = cfgInfo.MaxDistance;

            if (distance <= needDistance)
            {
                return (int)distance;
            }
        }
        return -1;
    }
}

