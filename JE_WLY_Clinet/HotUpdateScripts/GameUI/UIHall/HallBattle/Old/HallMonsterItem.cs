using HotUpdateScripts;
using JEngine.Helper;
using UnityEngine;

namespace HallBattle
{
    public class HallMonsterItem
    {
        private CFGItem_MainBattle cFGItem_MainBattle;
        public CFGItem_Monster cFGItem_Monster;
        public GameObject MonsterGo { get; private set; }
        public Transform MonsterTrans { get; private set; }
        public RectTransform MonsterRectTrans { get; private set; }
        public long HP { get; private set; }
        public long OriginHP { get; private set; }
        public int Hit { get; private set; } // 伤害
        public int SeckillTimes { get; private set; } // x次后直接秒杀英雄

        public float Scale { get; set; }

        private int totalMonsterCount;
        private int monsterIndex; // 怪物索引，每关可能多个怪物

        private bool useSkill;
        private int skillIndex = 0;

        public void OnInitConfig(CFGItem_MainBattle cfg)
        {
            cFGItem_MainBattle = cfg;
            monsterIndex = 1;
        }

        // 获取模型
        public string GetModelName()
        {
            if (cFGItem_MainBattle == null || cFGItem_MainBattle.monster == null || cFGItem_MainBattle.monster.Length < monsterIndex)
                return string.Empty;

            var info = cFGItem_MainBattle.monster[monsterIndex - 1];
            totalMonsterCount = cFGItem_MainBattle.monster.Length;
            HP = info[1]; // 血量
            OriginHP = HP;
            Hit = info[2]; // 攻击
            SeckillTimes = info[3]; // x次后秒杀英雄

            skillIndex = 0;
            cFGItem_Monster = CFG_Monster.GetById(info[0]);
            useSkill = (cFGItem_Monster.Skills != null && cFGItem_Monster.Skills.Length > 1);

            return cFGItem_Monster.Model;
        }

        public void OnInit(GameObject go)
        {
            MonsterGo = go;
            MonsterTrans = go.transform;
            MonsterRectTrans = go.GetComponent<RectTransform>();
        }

        public void Reset(Vector3 monsterPos, float screenWidth)
        {
            MonsterTrans.localPosition = monsterPos + new Vector3(screenWidth  / Scale, 0, 0);
            MonsterTrans.localScale = Vector3.one;
        }

        public bool IsRomate()
        {
            return CFG_Monster.IsRomate(cFGItem_Monster.ID);
        }

        public void DoAttackSpine(float playTimeScale)
        {
            var animName = "attack";

            if (useSkill && skillIndex != 0)
            {
                animName = "skill";
                skillIndex = 0;
            }
            else
                skillIndex = 1;

            SpineUtil.Instance.PlayOnce(MonsterTrans, animName, () =>
            {
                SpineUtil.Instance.PlayLoop(MonsterTrans, GameConst.ModelAnim.stand.ToString(), playTimeScale);
            }, playTimeScale);
        }

        #region position
        public Vector3 GetHpPostion()
        {
            return new Vector3(0, 100, 0) + MonsterTrans.parent.localPosition;
        }

        public Vector3 GetCenterPositon(int posType = 1)
        {
            var pos = MonsterTrans.position;
            pos.y += 0.6f * posType;
            return pos;
        }

        public Vector3 GetCenterLocalPositon()
        {
            var pos = MonsterTrans.localPosition;
            pos.x = 40;
            pos.y += 250 * MonsterTrans.parent.localScale.x;

            return pos;
        }
        #endregion

        public bool BeHit(long hit)
        {
            HP -= hit;
            return HP <= 0;
        }

        public void MoveToNext()
        {
            monsterIndex++;
            SetActive(false);
        }

        public void Clear()
        {
            if (MonsterGo != null)
            {
                MonsterGo.SetActive(false);
                GameObject.Destroy(MonsterGo);
            }

            MonsterGo = null;
            MonsterTrans = null;
        }

        public bool FinishLevel()
        {
            return monsterIndex >= totalMonsterCount;
        }

        public void SetActive(bool active)
        {
            if (MonsterGo != null)
            {
                MonsterGo.SetActive(active);

                if(!active)
                    MonsterTrans.localPosition =new Vector3(Screen.width * 2f, 0, 0);
            }
        }
    }
}