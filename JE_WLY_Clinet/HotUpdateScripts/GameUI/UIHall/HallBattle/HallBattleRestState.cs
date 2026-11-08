using HotUpdateScripts;
using System.Collections;
using UnityEngine;

namespace HallBattle
{
    public class HallBattleRestState : BaseState<UIHall>
    {
        public HallBattleRestState() : base(0.8f) { } // 设置延时

        private UIViewHall uIViewHall;

        private Vector3[] restPoss = new Vector3[5] { new Vector3(-162, 40, 0),  new Vector3(14, 186, 0), new Vector3(-110, 150, 0), new Vector3(146, 150, 0), new Vector3(158, 44, 0) };

        private bool inExit = false;

        public override void Enter()
        {
            uIViewHall = Owner.GetView();

            inExit = false;

            if(Owner.hallRole != null)
            {
                // 重置连击
                Owner.ResetCombo();
                Owner.SetMainLevelUIActive();
                Owner.SetMainBossHpUIActive();
                Owner.HideSpecialSkill();
            }

            ChangeBg();
            ChangeRoleStandPosition();

            AudioManager.Instance.PlayBgMusic(AudioKey.Sound.BG_HallRest);
        }

        private void ChangeBg()
        {
            uIViewHall.bgFarNodeTrans.gameObject.SetActive(false);
            uIViewHall.bgMidNodeTrans.gameObject.SetActive(false);
            uIViewHall.bgNearbyNodeTrans.gameObject.SetActive(false);
            uIViewHall.restBgGo.SetActive(true);
        }

        private void ChangeRoleStandPosition()
        {
            if (Owner.hallRole != null)
                Owner.hallRole.ChangeRoleState(HallRoleState.Rest);

            for (int i = 0; i < uIViewHall.roleTranss.Length; i++)
            {
                uIViewHall.roleTranss[i].localPosition = restPoss[i];
                uIViewHall.roleTranss[i].gameObject.SetActive(true);
            }
        }

        public override void Update(float dTime)
        {
            if (!inExit) return;
        }

        public override IEnumerator Exit()
        {
            inExit = true;

            Owner.hallRole.ChangeRoleState(HallRoleState.RestToBattle);
            Owner.ShowZhuanLoading();

            yield return base.Exit();
        }
    }
}