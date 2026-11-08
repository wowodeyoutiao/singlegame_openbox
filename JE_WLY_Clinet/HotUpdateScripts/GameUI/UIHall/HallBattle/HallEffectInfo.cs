using DG.Tweening;
using JEngine.Helper;
using System.Collections.Generic;
using UnityEngine;

namespace HallBattle
{
    public enum HallEffectColliderType
    {
        Box = 0,
        Circle = 1,
    }

    public enum HallEffectType
    {
        None = 0,
        Role,
        Monster
    }

    public enum HallHPChangeType
    {
        None = 0,
        Role = 1,
        FS = 2,
        WG = 3
    }

    public class HallEffectInfo : MonoBehaviour
    {
        [HideInInspector]
        private BoxCollider2D boxCollider;
        private CircleCollider2D circleCollider;

        [HideInInspector]
        public float AttackPower { get; private set; }
        [HideInInspector]
        public int AttackNum { get; private set; }
        [HideInInspector]
        public HallEffectType hallEffectType { get; private set; }
        [HideInInspector]
        public HallHPChangeType hallHPChangeType { get; private set; }
        [HideInInspector]
        public string HitEffectName { get; private set; }
        public Vector3 HitEffectPos { get; private set; }
        private int[] maxSize;
        [HideInInspector]
        public int skillMainCfgID;

        [HideInInspector]
        public int hitNum;
        [HideInInspector]
        public float hitTime;

        private Dictionary<string, int> woundedDict = new Dictionary<string, int>();

        public void OnInit(float power, int num, HallEffectType eType, int job = -1)
        {
            AttackPower = power;
            if (AttackPower < 1) AttackPower = 1;
            AttackNum = num;
            hallEffectType = eType;

            if (job == -1)
                hallHPChangeType = HallHPChangeType.Role;

            HitEffectName = string.Empty;
            woundedDict.Clear();
        }

        public void SetHitInfo(int num, float time)
        {
            hitNum = num;
            hitTime = time;
        }

        public void SetHPChangeType(bool isSpecialKill, int cfgID)
        {
            if (isSpecialKill)
                hallHPChangeType = HallHPChangeType.WG;
            else
                hallHPChangeType = HallHPChangeType.FS;

            skillMainCfgID = cfgID;
        }

        public void SetHitEffectInfo(string name, int[] posInfo)
        {
            if (string.IsNullOrEmpty(name) || posInfo == null || posInfo.Length < 2)
                return;

            HitEffectName = name;
            HitEffectPos = new Vector3(posInfo[0], posInfo[1], 0);
        }

        private void LoadCollider2D(HallEffectColliderType cType = HallEffectColliderType.Box)
        {
            if(cType == HallEffectColliderType.Circle)
            {
                if (boxCollider != null)
                {
                    boxCollider.enabled = false;
                    GameObject.DestroyImmediate(boxCollider);
                }

                if (circleCollider != null)
                    return;

                circleCollider = GetComponent<CircleCollider2D>();
                if (circleCollider == null)
                    circleCollider = gameObject.AddComponent<CircleCollider2D>();
            }
            else
            {
                if (circleCollider != null)
                {
                    circleCollider.enabled = false;
                    GameObject.DestroyImmediate(circleCollider);
                }

                if (boxCollider != null)
                    return;

                boxCollider = GetComponent<BoxCollider2D>();
                if (boxCollider == null)
                    boxCollider = gameObject.AddComponent<BoxCollider2D>();
            }
        }

        public void SetSize(int[] size, HallEffectColliderType cType = HallEffectColliderType.Box)
        {
            LoadCollider2D(cType);

            if (cType == HallEffectColliderType.Circle)
                circleCollider.radius = size[0];
            else 
                boxCollider.size = new Vector2(size[0], size[1]);

            maxSize = size;
        }

        public void SetSize(float value, HallEffectColliderType cType = HallEffectColliderType.Box)
        {
            LoadCollider2D(cType);

            if (cType == HallEffectColliderType.Circle)
                circleCollider.radius = value;
            else
                boxCollider.size = new Vector2(maxSize[0], maxSize[1]) * value;
        }

        public void SetOffset(int[] offset, HallEffectColliderType cType = HallEffectColliderType.Box)
        {
            if (offset == null || offset.Length == 0)
                return;

            LoadCollider2D(cType);

            if (cType == HallEffectColliderType.Circle)
                circleCollider.offset = new Vector2(offset[0], offset[1]);
            else
                boxCollider.offset = new Vector2(offset[0], offset[1]);
        }

        public void ChangeAttackNum(int value)
        {
            AttackNum += value;
        }

        public void AddWounded(string name)
        {
            if(!woundedDict.TryGetValue(name, out var index))
                woundedDict.Add(name, index);
        }

        public bool CanHit(string name)
        {
            return !woundedDict.ContainsKey(name);
        }

        public void EnableCollider()
        {
            if (boxCollider != null)
                boxCollider.enabled = true;

            if (circleCollider != null)
                circleCollider.enabled = true;
        }

        public void DisableCollider()
        {
            if (boxCollider != null)
                boxCollider.enabled = false;

            if (circleCollider != null)
                circleCollider.enabled = false;
        }

        public void SizeZoomTo(float value, float time)
        {
            DOTween.To(() => 1, (a) =>
            {
                if (maxSize.Length == 1)
                    circleCollider.radius = a * 0.1f * maxSize[0];
                else
                    boxCollider.size = new Vector2(maxSize[0], maxSize[1]) * a * 0.1f;

            }, 10, time).OnComplete(() =>
            {
                DisableCollider();
                PoolHelper.Instance.RecycObj(gameObject);
            });
        }
    }
}