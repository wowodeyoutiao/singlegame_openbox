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

public class Monster : MonoBehaviour
{
    public SkeletonAnimation anim;
    public GameConst.ModelAnim _animState;
    public bool isEndMove = false;

    //public NavMeshAgent agent;
    public int _priority;
    public Point _gridNow;//当前格子

    public double _gridX, _gridY;

    private float _maxLifeTime;//最大时间
    public float maxLifeTime
    {
        get => _maxLifeTime;

        set
        {
            if (_maxLifeTime == value)
                return;
            _maxLifeTime = value;
        }
    }

    private float _curLifeTime;//当前时间
    public float curLifeTime
    {
        get => _curLifeTime;

        set
        {
            if (_curLifeTime == value)
                return;
            _curLifeTime = value;
        }
    }

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

    private string _monsterName;//怪物名字
    public string monsterName
    {
        get => _monsterName;

        set
        {
            if (_monsterName.CompareTo(value) == 0)
                return;
            _monsterName = value;

        }
    }

    private long _uID;//ID
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

    private int _cfgId;//ID
    public int cfgId
    {
        get => _cfgId;

        set
        {
            if (_cfgId == value)
                return;
            _cfgId = value;

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

    /// <summary>
    /// 怪物的归属
    /// </summary>
    private Ownership _ownership;

    /// <summary>
    /// 怪物的归属
    /// </summary>
    public Ownership _Ownership
    {
        get => _ownership;

        set
        {
            if (_ownership == value)
                return;
            _ownership = value;

        }
    }

    private string BulletIcon = "yg_01";

    private Dictionary<string, GameObject> GameObjDic = new Dictionary<string, GameObject>();
    public List<int> _BufferList = new List<int>();

    private void AddGameObj(string key, GameObject jac)
    {
        if (!GameObjDic.ContainsKey(key))
        {
            GameObjDic.Add(key, jac);
        }
        else
        {
            ClearOneObj(key);
            GameObjDic[key] = null;
            GameObjDic[key] = jac;
        }
    }

    private void ClearOneObj(string key)
    {
        if (GameObjDic.ContainsKey(key))
        {
            PoolHelper.Instance.RecycObj(GameObjDic[key]);
            GameObjDic.Remove(key);
        }
    }

    private void CleacGameObjDic()
    {
        foreach (var item in GameObjDic.Values)
        {
            PoolHelper.Instance.RecycObj(item);
        }
        GameObjDic.Clear();
    }

    private void OnDisable()
    {
        CleacGameObjDic();
    }

    public async UniTaskVoid AddOrRemoveBuffer(int bufferId, bool isAdd = true)
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

        if (GameRoomM.Instance.GameAttack == GameRoomM.AttackClass.Team && bufferId != 301 && bufferId != 401 && bufferId != 501 && bufferId != 601)
            return;

        string effKey = CFG_BuffEffect.GetById(bufferId).BuffEffect;

        if (!string.IsNullOrEmpty(effKey) && GameConst.CheakEff(effKey))
        {
            if (isAdd)
            {
                GameObject effObj = null;
                string poolKey = effKey;

                if (!PoolHelper.Instance.CheckPool(poolKey))
                {
                    var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
                    await UniTask.WaitUntil(() => baseObj.Loaded);
                    PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
                }
                await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
                effObj = PoolHelper.Instance.GetPoolObj(poolKey);
                if (gameObject == null || gameObject.Equals(null))
                {
                    PoolHelper.Instance.RecycObj(effObj);
                    return;
                }
                effObj.transform.SetParent(this.transform);
                effObj.transform.localPosition = Vector3.back * 0.1f;
                effObj.SetActive(true);
                AddGameObj(effKey, effObj);
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
                AddOrRemoveBuffer(item.BuffId, true).Forget();
                _BufferList.Add(item.BuffId);
                Idx++;
            }
        }

    }

    public Transform root_hub;
    public Transform root_eff;
    public CFGItem_Monster _MonsterCFG;

    public void OnInit()
    {
        anim = transform.Find("spine").GetComponent<SkeletonAnimation>();
        root_hub = transform.Find("root_hub");

        root_eff = null;
        if (transform.Find("root_Eff"))
            root_eff = transform.Find("root_Eff");

        if (_gridNow == null)
            _gridNow = new Point();
        //agent.avoidancePriority = _priority;
        //agent.speed = 2;
        //agent.radius = 0.1f;

        _animState = GameConst.ModelAnim.No;
        changeState(GameConst.ModelAnim.stand);

        _Dir = Direction.LEFT;

        _MonsterCFG = CFG_Monster.GetById(cfgId);
        if (_MonsterCFG.Isboss == 3)
            maxLifeTime = _MonsterCFG.Delay / 1000f;

        curLifeTime = 0;

    }


    public void changeState(GameConst.ModelAnim actName,float AnimAttackTime = 3f)
    {
        // Debug.Log("actName " + actName + "  _animState  "+ _animState);
        if (actName != GameConst.ModelAnim.dead && actName != GameConst.ModelAnim.attack && !_animState.ToString().Contains("skill") && _animState.CompareTo(actName) == 0)
            return;
        _animState = actName;

        if (_animState == GameConst.ModelAnim.attack || _animState.ToString().Contains("skill"))
        {
            SpineUtil.Instance.PlayOnce(anim, actName.ToString(), () =>
            {
                _animState = GameConst.ModelAnim.stand;
                SpineUtil.Instance.PlayLoop(anim, _animState.ToString());
            }, AnimAttackTime);
        }

        else if (actName == GameConst.ModelAnim.dead)
        {
            //箱子，袋子开启声音
            if ((this._cfgId >= 25001 && this._cfgId <= 25148) || (this.cfgId >= 602 && this.cfgId <= 607))
                AudioManager.Instance.PlaySound(AudioKey.Sound.UI_OpenBox);
            else
                AudioManager.Instance.PlaySound(AudioKey.Sound.Battle_MonsterDied);

            SpineUtil.Instance.PlayOnce(anim, actName.ToString(), null);
        }
        else
            SpineUtil.Instance.PlayLoop(anim, actName.ToString());
    }

    private HashSet<string> EffHasSet = new HashSet<string>();

    public async UniTaskVoid SetHurtChange(int changeNum, long AttackerId, int skillId)
    {
        if (GameRoomM.Instance.GameAttack == GameRoomM.AttackClass.Team && AttackerId != MainRole._MainRole.uID)
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

            //AudioManager.Instance.PlaySound((AudioKey.Sound)Enum.Parse(typeof(AudioKey.Sound), poolKey));
            AudioManager.Instance.PlaySound((AudioKey.Sound)Enum.Parse(typeof(AudioKey.Sound), "Battle_MonsterHit"));

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

            EffHasSet.Add(poolKey);
            CoroutineCall(0.8f, () =>
            {
                EffHasSet.Remove(poolKey);
                PoolHelper.Instance.RecycObj(Eff);
            });
        }
    }

    public async UniTask SetAttack(Point Location, Direction direction, int SkillId, long[] TargetUids)
    {
        if (_animState == GameConst.ModelAnim.dead)
            return;

        setRotation(direction);

        var CFGSKill = CFG_Skill.GetById(SkillId);
        changeState(GameConst.GetAnim((int)CFGSKill.Needskill), _MonsterCFG.Isboss == 4? 1 : 3f);

        if (GameRoomM.Instance.GameAttack == GameRoomM.AttackClass.Team
            && TargetUids != null && !GameRoomM.Instance.CheakIsMainRole(TargetUids[0]))
            return;

        var skillEff = CFG_SkillEffect.GetById(CFGSKill.SkillID);

        var baseRoot = root_eff != null ? root_eff : transform;

        if (skillEff != null)
        {
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
                    Eff.transform.position = new Vector3(baseRoot.position.x, baseRoot.position.y, SceneManager.EffOffZ);

                    if (skillEff.isAmendStart != 0)
                        Eff.transform.eulerAngles = SceneManager.GetSeverRotationPos(direction);
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

            if (Location != null)
            {
                var tpos = new Vector3((float)Location.X, (float)Location.Y, SceneManager.GetOffZ(Location.X, Location.Y));
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
                    Eff.transform.position = new Vector3(baseRoot.position.x, baseRoot.position.y, SceneManager.EffOffZ);

                    Vector2 _direction = (baseRoot.position - tpos);
                    _direction.Normalize();
                    // 设置物体的朝向
                    Eff.transform.rotation = Quaternion.AngleAxis(Mathf.Atan2(_direction.y, _direction.x) * Mathf.Rad2Deg, Vector3.forward);

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
        }
        else if (Location != null)
        {
            string poolKey = _MonsterCFG.AttackEffect;

            var SoundKey = AudioKey.Sound.BossAttack;
            if (_MonsterCFG.Isboss == 0 && !string.IsNullOrEmpty(poolKey))
            {
                if (poolKey.Contains("yg_"))
                {
                    poolKey = "Monster_zd";
                    BulletIcon = _MonsterCFG.AttackEffect;

                    SoundKey = string.IsNullOrEmpty(poolKey) ? AudioKey.Sound.JZAttack : AudioKey.Sound.YCAttack;
                }
            }

            bool isFly = true;
            if (_MonsterCFG.Model == "Boss_Romate_1")
            {
                isFly = false;
                poolKey = "Boss_Romate_1_Bullet";
            }
            else if (_MonsterCFG.Model == "Boss_Romate_2")
            {
                isFly = false;
                poolKey = "Boss_Romate_2_Bullet";
            }

            AudioManager.Instance.PlaySound(SoundKey);

            if (string.IsNullOrEmpty(poolKey))
                return;
            var tpos = new Vector3((float)Location.X, (float)Location.Y, SceneManager.GetOffZ(Location.X, Location.Y));
            if (!PoolHelper.Instance.CheckPool(poolKey))
            {
                var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
                await UniTask.WaitUntil(() => baseObj.Loaded);
                PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
            }
            await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(poolKey));
            var Eff = PoolHelper.Instance.GetPoolObj(poolKey);
            Eff.transform.SetParent(FightController.Instance.Root_Eff3D);
            Eff.SetActive(true);

            if (isFly)
            {
                Eff.transform.position = new Vector3(transform.position.x, transform.position.y, SceneManager.EffOffZ);
                if (poolKey == "Monster_zd")
                {
                    ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.FightAtlas, BulletIcon, Eff.GetComponent<SpriteRenderer>()).Forget();
                }
                Vector2 _direction = (transform.position - tpos);
                _direction.Normalize();
                Eff.transform.rotation = Quaternion.AngleAxis(Mathf.Atan2(_direction.y, _direction.x) * Mathf.Rad2Deg, Vector3.forward);

                Eff.transform.DOMove(tpos, 0.15f).OnComplete(() =>
                {
                    Eff.transform.DOKill();
                    PoolHelper.Instance.RecycObj(Eff);
                });
            }
            else
            {
                Eff.transform.position = tpos;
                CoroutineCall(0.3f, () =>
                {
                    if (Eff != null)
                        PoolHelper.Instance.RecycObj(Eff);
                });
            }

        }
    }


    private void Update()
    {
        //if (isEndMove)
        //{
        //    // 计算当前帧的移动速度
        //    if (Vector3.Distance(transform.position, _runToPos) < 0.01f)
        //    {
        //        isEndMove = false;
        //        changeState(GameConst.ModelAnim.stand);
        //    }
        //    else
        //    {
        //        // 执行移动
        //        transform.position = Vector3.MoveTowards(
        //            transform.position,
        //            _runToPos,
        //            requiredSpeed * UnityEngine.Time.deltaTime
        //        );
        //    }
        //}
    }
    private float requiredSpeed = 0;
    private Vector3 _runToPos;
    private float _runToTime = 0.2f; //移动时间

    public void moveTo(Point Location, Direction direction)
    {
        if (_animState == GameConst.ModelAnim.dead)
            return;

        float offZ = SceneManager.GetOffZ(Math.Floor(Location.X), Math.Floor(Location.Y));
        if (offZ != SceneManager.Block)
        {
            changeState(GameConst.ModelAnim.walk);
            isEndMove = true;
            _runToPos = new Vector3((float)Location.X, (float)Location.Y, offZ);
            _runToTime = 0.8f;
            requiredSpeed = Vector3.Distance(transform.position, _runToPos) / _runToTime;
            //agent.SetDestination(_runToPos);
            this.transform.DOKill();
            this.transform.DOMove(_runToPos, 0.4f).SetEase(Ease.Linear).OnComplete(() =>
            {
                isEndMove = false;
                changeState(GameConst.ModelAnim.stand);
            });
            setRotation(direction);

            _gridNow.X = Math.Floor(Location.X);
            _gridNow.Y = Math.Floor(Location.Y);

            _gridX = _gridNow.X;
            _gridY = _gridNow.Y;
        }

    }

    //传送
    public async UniTaskVoid SetTeleport(Point point, Direction dir, TeleportType Type)
    {
        wrap(new Vector3((float)point.X, (float)point.Y, SceneManager.GetOffZ(Math.Floor(point.X), Math.Floor(point.Y))));
        setRotation(dir);

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
            Eff.transform.position = this.transform.position;
            Eff.SetActive(true);
            CoroutineCall(1, () =>
            {
                if (Eff != null)
                    PoolHelper.Instance.RecycObj(Eff);
            });
        }
    }

    public void wrap(Vector3 target, bool isInit = false)
    {
        //agent.Warp(target);
        this.transform.DOKill();
        this.transform.position = target;

        _gridNow.X = target.x;
        _gridNow.Y = target.y;

        _gridX = _gridNow.X;
        _gridY = _gridNow.Y;

        //if (isInit)
        //{
        //string poolKey = "eff_MonsterBorn";
        //if (!PoolHelper.Instance.CheckPool(poolKey))
        //{
        //    var baseObj = new JPrefab(GameConst.GetEffPath(poolKey), true);
        //    await baseObj.WaitForAsyncLoading();
        //    PoolHelper.Instance.AddPool(poolKey, baseObj.Instance);
        //}

        //var Eff = PoolHelper.Instance.GetPoolObj(poolKey);
        //Eff.transform.SetParent(FightController.Instance.Root_Eff3D);
        //Eff.transform.position = this.transform.position;
        //Eff.SetActive(true);
        //CoroutineCall(1.5f, () =>
        //{
        //    PoolHelper.Instance.RecycObj(Eff);
        //});
        //}
    }

    public void setRotation(Direction dir)
    {
        if (dir != Direction.UP && dir != Direction.DOWN)
        {
            _Dir = dir;
        }

        Vector3 endPos = Vector3.one;
        if (_Dir == Direction.UP_LEFT || _Dir == Direction.DOWN_LEFT || _Dir == Direction.LEFT)
            endPos = Vector3.zero;
        else if (_Dir == Direction.UP_RIGHT || _Dir == Direction.DOWN_RIGHT || _Dir == Direction.RIGHT)
            endPos = new Vector3(0, 180, 0);

        this.transform.eulerAngles = endPos;
    }


    public bool IsMoving()
    {
        return isEndMove;
    }

    // private JAction DesJa;
    public void destoryMonster()
    {
        isEndMove = false;
        _MonsterCFG = null;
        changeState(GameConst.ModelAnim.dead);
        transform.DOKill();

        CoroutineCall(1.2f, () =>
        {
            if (this.gameObject != null && !this.gameObject.Equals(null))
            {
                ClearCoroutine();
                CleacGameObjDic();
                string pooKey = this.gameObject.name.Split('/')[0];
                PoolHelper.Instance.RecycObj(pooKey, this.gameObject);
            }
        });
    }

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

    private System.Collections.IEnumerator IECoroutineCall(float time, Action callBack, int coroutineId)
    {
        yield return new WaitForSeconds(time);
        coroutinesDic.Remove(coroutineId);
        callBack?.Invoke();
    }

    public void Destory()
    {
        ClearCoroutine();
        CleacGameObjDic();
        //if (DesJa != null)
        //{
        //    DesJa.Cancel();
        //}
        //PoolHelper.Instance.RecycObj("Monser" + cfgId, this.gameObject);
    }

    void OnDestroy()
    {
        Destory();
    }

    /// <summary>
    /// 获取技能伤害距离
    /// </summary>
    /// <param name="skillID">为0时 按照dis 距离检测 ，其他则按照技能归宿角色ID</param>
    /// <param name="dis"></param>
    /// <returns></returns>
    public int getHurtDistance(int skillID, int dis, bool isAttack)
    {
        if (isAttack && hp <= 0)
            return -1;

        if (!isAttack && maxLifeTime <= 0)
            return -1;

        if (_animState == GameConst.ModelAnim.dead)
            return -1;

        // int distance = MainRole._MainRole.getDistanceFromMainRole(_gridNow.X, _gridNow.Y);
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

