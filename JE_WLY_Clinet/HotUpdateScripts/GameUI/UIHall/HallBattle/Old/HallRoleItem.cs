using HotUpdateScripts.ProtoM;
using UnityEngine;

namespace HallBattle
{
    public class HallRoleItem
    {
        public Transform RoleTrans { get; set; }
        public bool IsMainRole { get; private set; } // 主角
        public int ConfigId { get; private set; } // 角色表ID
        public long Uid { get; private set; }
        public int SoltId { get; private set; }
        public int Job { get; private set; }

        private long hP;
        private long curHP;

        public void OnInt(bool isMainRole, int job)
        {
            IsMainRole = isMainRole;
            Job = job;
        }

        public void SetPartnerInfo(int configId, long uid, int soltId)
        {
            ConfigId = configId;
            Uid = uid;
            SoltId = soltId;

            ResetHP();
        }

        public Vector3 GetCenterPosition()
        {
            return RoleTrans.localPosition + RoleTrans.parent.localPosition + RoleTrans.parent.parent.localPosition;
        }

        public Vector3 GetInsideCenterPosition()
        {
            return RoleTrans.localPosition + RoleTrans.parent.localPosition;
        }

        public void ResetHP()
        {
            if (IsMainRole)
                hP = RoleM.Instance.GetRoleAttrValueByAttr(pb.Attribute.MAX_HEALTH);
            else
                hP = PartnerM.Instance.GetRoleAttrValueByAttr(Uid, pb.Attribute.MAX_HEALTH);

            curHP = hP;
        }

        public int GetFullHP()
        {
            return (int)hP;
        }

        public int GetHit()
        {
            if (IsMainRole)
                return RoleM.Instance.Combat;
            else
            {
                var partnerStats = PartnerM.Instance.GetOnePartnerStats(Uid);
                if (partnerStats != null)
                    return partnerStats.Combat;
            }
            return 0;
        }

        public bool BeHit(long value)
        {
            curHP -= value;
            return curHP <= 0;
        }

        public bool IsDie()
        {
            return curHP <= 0;
        }
    }
}