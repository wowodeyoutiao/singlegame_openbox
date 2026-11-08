using Cysharp.Threading.Tasks;
using DG.Tweening;
using HotUpdateScripts;
using JEngine.Helper;
using UnityEngine;
using UnityEngine.UI;

namespace HallBattle
{
    public class HallMonster : HallBaseMonster
    {
        private string romateMonsterEffectKey = "Romate_Monster_UI_Effect";

        public override void Attack(HallRole target)
        {
            if (beHitTimes >= seckillTimes && !isMopUp)
            {
                //Debug.LogError($"结束战斗, {gameObject.name}，{beHitTimes}, {seckillTimes}");
                UniGameEvent.Send(new GameEvent.HallBattleSettlement() { IsRoleWin = false });
            }

            base.Attack(target);

            if (IsRomate)
                ShowRomateMonsterAttackEffect(target).Forget();
            else
            {
                // 发射
                FireBoxCollider2D(target.transform.position, 0.8f).Forget();
            }
        }

        private async UniTask FireBoxCollider2D(Vector3 position, float time)
        {
            if (curMonsterState != HallMonsterState.Attack)
                return;

            var effectKey = "Monster_UI_Effect_Info";
            if (!PoolHelper.Instance.CheckPool(effectKey))
            {
                var go = new GameObject(effectKey);
                PoolHelper.Instance.AddPool(effectKey, go);
            }
            await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(effectKey));
            var effGo = PoolHelper.Instance.GetPoolObj(effectKey);
            effGo.name = effectKey;
            effGo.transform.SetParent(transform.parent);
            effGo.transform.position = transform.position;
            effGo.transform.localScale = Vector3.one;
            effGo.SetActive(true);

            var hallEffectInfo = effGo.GetComponent<HallEffectInfo>();
            if (hallEffectInfo == null)
                hallEffectInfo = effGo.AddComponent<HallEffectInfo>();

            hallEffectInfo.OnInit(attackPower * (105 - Random.Range(0, 11)) / 100, 1, HallEffectType.Monster);
            hallEffectInfo.SetSize(new int[] { 50, 200 });

            position.y += 0.5f;
            position.z = 0;
            effGo.transform.DOMove(position, time).OnComplete(() =>
            {
                if (effGo != null)
                    PoolHelper.Instance.RecycObj(effGo);
            });
        }

        private async UniTask ShowRomateMonsterAttackEffect(HallRole target)
        {
            string attackEffect = cFGItem_Monster.AttackEffect;
            var SoundKey = AudioKey.Sound.BossAttack;
            string bulletIcon = "";

            if (cFGItem_Monster.Isboss == 0 && !string.IsNullOrEmpty(attackEffect))
            {
                if (attackEffect.Contains("yg_"))
                {
                    attackEffect = romateMonsterEffectKey;
                    bulletIcon = cFGItem_Monster.AttackEffect;

                    SoundKey = string.IsNullOrEmpty(attackEffect) ? AudioKey.Sound.JZAttack : AudioKey.Sound.YCAttack;
                }
            }

            AudioManager.Instance.PlaySound(SoundKey);

            if (string.IsNullOrEmpty(attackEffect) || attackEffect != romateMonsterEffectKey)
                return;

            if (!PoolHelper.Instance.CheckPool(attackEffect))
            {
                var go = new GameObject(romateMonsterEffectKey);
                go.AddComponent<Image>();
                PoolHelper.Instance.AddPool(attackEffect, go);
            }
            await UniTask.WaitUntil(() => PoolHelper.Instance.CheckPool(attackEffect));
            var effGo = PoolHelper.Instance.GetPoolObj(attackEffect);
            effGo.name = romateMonsterEffectKey;
            effGo.transform.SetParent(transform.parent);
            effGo.transform.position = transform.position + new Vector3(-2.5f * scale, 1.8f * scale);
            effGo.transform.localScale = Vector3.one;
            effGo.SetActive(false);

            var hallEffectInfo = effGo.GetComponent<HallEffectInfo>();
            if(hallEffectInfo == null)
                hallEffectInfo = effGo.AddComponent<HallEffectInfo>();

            hallEffectInfo.OnInit(attackPower * (105 - Random.Range(0, 11)) / 100, 1, HallEffectType.Monster);
            hallEffectInfo.SetSize(new int[] { 50, 50 });

            var img = effGo.GetComponent<Image>();
            ResouseHelper.Instance.LoadSpriteByAtlas(GameConst.FightAtlas, bulletIcon, img, callBack: () => { img.SetNativeSize();  });

            GameHelper.Instance.SetDeayTimeDo(() =>
            {
                if (effGo == null)
                    return;

                effGo.SetActive(true);

                var startPos = effGo.transform.position;
                var tpos = target.transform.position;
                tpos.y += 0.2f;

                Vector2 dir = (startPos - tpos);
                dir.Normalize();
                effGo.transform.rotation = Quaternion.AngleAxis(Mathf.Atan2(dir.y, dir.x) * Mathf.Rad2Deg, Vector3.forward);

                effGo.transform.DOMove(tpos, 0.2f).SetEase(Ease.Linear).OnComplete(() =>
                {
                    effGo.transform.DOKill();
                    PoolHelper.Instance.RecycObj(effGo);
                });
            }, 0.5f);
        }
    }
}