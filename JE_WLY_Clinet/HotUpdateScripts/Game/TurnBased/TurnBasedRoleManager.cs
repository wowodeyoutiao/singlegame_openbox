using JEngine.Singleton;

using pb;

using System.Collections.Generic;

using Object = UnityEngine.Object;

namespace HotUpdateScripts.Game
{
    public class TurnBasedRoleManager : SingletonTemplate<TurnBasedRoleManager>
    {
        public Dictionary<int, TurnBasedRole> blueRoleDic = new Dictionary<int, TurnBasedRole>();
        public Dictionary<int, TurnBasedRole> redRoleDic = new Dictionary<int, TurnBasedRole>();

        public List<TurnBasedRole> allRoles = new List<TurnBasedRole>();

        public void AddBuleRole(TurnBasedRole role, int position)
        {
            if (role == null) return;
            blueRoleDic[position] = role;
            allRoles.Add(role);
        }

        public void AddRedRole(TurnBasedRole role, int position)
        {
            if (role == null) return;
            redRoleDic[position] = role;
            allRoles.Add(role);
        }

        /*
        public bool GetBlueRoleByPosition(int position, out TurnBasedRole role)
        {
            bool hasRole = blueRoleDic.TryGetValue(position, out TurnBasedRole _role);
            if (!Object.Equals(hasRole, null))
            {
                role = _role;
                return hasRole;
            }
            role = null;
            return false;
        }

        public bool GetRedRoleByPosition(int position, out TurnBasedRole role)
        {
            bool hasRole = redRoleDic.TryGetValue(position, out TurnBasedRole _role);
            if(!Object.Equals(hasRole, null))
            {
                role = _role;
                return hasRole;
            }
            role = null;
            return false;
        }
        */

        public bool GetRoleByUid(long uid, out TurnBasedRole role)
        {
            var _role = allRoles.Find(r => r.uid == uid);
            if (!Equals(_role, null))
            {
                role = _role;
                return _role != null;
            }
            role = null;
            return false;
        }

        public bool GetBlueRoleByUid(long uid, out TurnBasedRole role)
        {
            foreach (var roleObj in blueRoleDic.Values)
            {
                if (!Equals(roleObj, null) && roleObj.uid == uid)
                {
                    role = roleObj;
                    return true;
                }
            }
            role = null;
            return false;
        }

        public bool GetRedRoleByUid(long uid, out TurnBasedRole role)
        {
            foreach (var roleObj in redRoleDic.Values)
            {
                if (!Equals(roleObj, null) && roleObj.uid == uid)
                {
                    role = roleObj;
                    return true;
                }
            }
            role = null;
            return false;
        }

        public TurnBasedBattle.ActionCamp GetRoleCampType(long uid)
        {
            var _role = allRoles.Find(r => r.uid == uid);
            if (!Equals(_role, null))
            {
                return _role.campType;
            }

            return TurnBasedBattle.ActionCamp.INVALID_CAMP;
        }

        public bool GetBlueProtagonist(out TurnBasedRole role)
        {
            foreach (var _role in blueRoleDic.Values)
            {
                if(_role.roleType == ObjectType.PLAYER)
                {
                    role = _role;
                    return true;
                }
            }
            role = null;
            return false;
        }

        public bool GetRedProtagonist(out TurnBasedRole role)
        {
            foreach (var _role in redRoleDic.Values)
            {
                if(_role.roleType == ObjectType.PLAYER)
                {
                    role = _role;
                    return true;
                }
            }
            role = null;
            return false;
        }

        public void ClearAllRole()
        {
            blueRoleDic.Clear();
            redRoleDic.Clear();

            foreach (var role in allRoles)
            {
                if (!Equals(role, null))
                {
                    Object.Destroy(role.gameObject);
                }
            }
            allRoles.Clear();
        }

        public void StopAllRoleAttack()
        {
            foreach (var role in allRoles)
            {
                if (!Equals(role, null))
                    role.StopAttack();
            }
        }

        public void OnDestroy()
        {
            StopAllRoleAttack();
            ClearAllRole();
        }
    }
}
