using HotUpdateScripts;
using JEngine.Helper;
using Spine.Unity;
using System.Collections.Generic;
using UnityEngine;

namespace HallBattle
{
    public class Combatant : MonoBehaviour
    {
        [Header("基本属性")]
        public long maxHealth = 100; // 初始血量
        public long currentHealth; // 当前血量
        public int attackPower = 10; // 战力

        protected Vector2 targetPosition; // 战斗位置

        // 移速、攻击范围
        public float moveSpeed = 1f;
        public float attackRange = 1.5f;

        public float timeScale = 1; // 攻速
        public float attackCooldown = 1f;// 攻击冷却
        protected float lastAttackTime; // 上次攻击时间

        [Header("战斗状态")]
        public bool isAlive = true; // 活着
        public bool isMoving = false; // 正在移动
        public bool isAttacking = false; // 正在攻击
        public bool isTakingDamage = false; // 正在受击

        [Header("引用")]
        public Rigidbody2D rb;
        public Transform attackPoint; // 攻击位置

        [Header("特效")]
        public GameObject hitEffect; // 受击特效
        public GameObject attackEffect; // 攻击特效
        public GameObject deathEffect; // 死亡特效

        private string spineName;
        protected Transform spineTrans;
        protected List<Combatant> targets = new List<Combatant>(); // 攻击目标

        protected virtual void Awake()
        {
            if (rb == null)
            {
                rb = gameObject.AddComponent<Rigidbody2D>();
                rb.bodyType = RigidbodyType2D.Kinematic;

                var collider2D = gameObject.AddComponent<BoxCollider2D>();
                collider2D.isTrigger = true;

                collider2D.offset = new Vector2(0, 80);
                collider2D.size = new Vector2(160, 300);
            }

            if (spineTrans == null)
                spineTrans = GetComponentInChildren<SkeletonGraphic>().transform;
        }

        public virtual void Birth(Vector2 birthPosition)
        {
            transform.localPosition = birthPosition;
            if (rb != null)
                rb.velocity = Vector2.zero;
        }

        protected virtual void Update()
        {
            if (!isAlive || isTakingDamage) return;

            var target = GetNearestTarget();
            if(target != null && IsInAttackRange(target))
                TryAttack();
            else if(isMoving)
                DoMove(target);
        }

        // 获取最近目标
        protected Combatant GetNearestTarget()
        {
            if (targets == null || targets.Count == 0)
                return null;

            Combatant target = null;
            var minDis = 9999f;
            for(int i = 0; i < targets.Count; i++)
            {
                if (!targets[i].isAlive)
                    continue;

                var dis = Vector2.Distance(transform.localPosition, targets[i].transform.localPosition);

                if(dis < minDis)
                {
                    target = targets[i];
                    minDis = dis;
                }
            }

            return target;
        }

        public virtual void GoToRest()
        {
            targets.Clear();
            isMoving = false;
            PlayLoopSpine(GameConst.ModelAnim.stand.ToString());
            Birth(Vector2.zero);
        }

        #region 移动
        protected virtual void DoMove(Combatant target)
        {
            Vector2 disV2 = Vector2.zero;
            if (target == null)
                disV2 = targetPosition - (Vector2)transform.localPosition;
            else
                disV2 = (Vector2)target.transform.localPosition - (Vector2)transform.localPosition;

            disV2.y = 0; // 只有X移动

            if(Mathf.Abs(disV2.x) < 10)
                rb.velocity = Vector2.zero;
            else
            {
                PlayLoopSpine(GameConst.ModelAnim.run.ToString());

                rb.velocity = disV2.normalized * moveSpeed;
            }  
        }

        public virtual void MoveTo(Vector2 pos, float moveSpeed)
        {
            targetPosition = pos;
            this.moveSpeed = moveSpeed;
            isMoving = true;
        }
        #endregion

        #region 攻击
        // 检查是否在攻击范围内
        protected virtual bool IsInAttackRange(Combatant target)
        {
            if (target == null) return false;
            return Vector2.Distance(transform.localPosition, target.transform.localPosition) <= attackRange;
        }

        // 尝试攻击
        protected virtual void TryAttack()
        {
            if (Time.time - lastAttackTime >= attackCooldown)
            {
                Attack();
            }
        }

        // 攻击逻辑
        public virtual void Attack()
        {
            isMoving = false;
            isAttacking = true;
            PlayLoopSpine(GameConst.ModelAnim.stand.ToString());
            lastAttackTime = Time.time;

            //// 播放攻击特效
            //if (attackEffect != null && attackPoint != null)
            //{
            //    Instantiate(attackEffect, attackPoint.position, Quaternion.identity);
            //}

            // 对目标造成伤害
            //if (target != null && IsInAttackRange())
            //{
            //    target.TakeDamage(attackPower);
            //}

            Invoke(nameof(EndAttack), 0.5f);
        }

        protected virtual void EndAttack()
        {
            isAttacking = false;
        }
        #endregion

        #region 受击
        // 受击逻辑
        public virtual void TakeDamage(int damage)
        {
            currentHealth = currentHealth - damage;
            if (currentHealth < 0) currentHealth = 0;

            isTakingDamage = true;
            //animator.SetTrigger("TakeDamage");

            // 播放受击特效
            if (hitEffect != null)
            {
                Instantiate(hitEffect, transform.position, Quaternion.identity);
            }

            // 击退效果
            //if (target != null)
            //{
            //    Vector2 knockbackDirection = ((Vector2)transform.position - (Vector2)target.transform.position).normalized;
            //    rb.velocity = knockbackDirection * 3f;
            //}

            if (currentHealth <= 0)
                Die();
            else
                Invoke(nameof(EndTakingDamage), 0.3f);
        }

        // 结束受击状态
        protected virtual void EndTakingDamage()
        {
            isTakingDamage = false;
            rb.velocity = Vector2.zero;
        }
        #endregion

        #region 死亡
        // 死亡逻辑
        public virtual void Die()
        {
            isAlive = false;
            isMoving = false;
            // animator.SetTrigger("Die");
            rb.velocity = Vector2.zero;

            // 播放死亡特效
            if (deathEffect != null)
            {
                Instantiate(deathEffect, transform.position, Quaternion.identity);
            }

            targets.Clear();
            // 移除目标引用
            //if (target != null)
            //{
            //    Combatant targetCombatant = target;
            //    target = null;
            //    if (targetCombatant.target == this)
            //    {
            //        targetCombatant.target = null;
            //    }
            //}

            // Destroy(gameObject, 1f);
        }
        #endregion

        public virtual void PlayLoopSpine(string name)
        {
            if (spineName == name)
                return;

            spineName = name;
            SpineUtil.Instance.PlayLoop(spineTrans, name, timeScale);
        }

        //protected virtual void Update()
        //{
        //    if (!isAlive || isTakingDamage) return;

        //    if (target != null && IsInAttackRange())
        //    {
        //        TryAttack();
        //    }
        //    else if (isMoving)
        //    {
        //        MoveToTargetPosition();
        //    }
        //    else
        //        DoLoopSpine(GameConst.ModelAnim.run.ToString(), timeScale);
        //}

        //// 移动到目标位置
        //public virtual void MoveTo(Vector2 position)
        //{
        //    targetPosition = position;
        //    isMoving = true;
        //    DoLoopSpine(GameConst.ModelAnim.run.ToString(), timeScale);
        //}



        //// 翻转角色朝向
        //protected virtual void Flip(float directionX)
        //{
        //    if (directionX != 0)
        //        transform.localScale = new Vector3(Mathf.Sign(directionX), 1, 1);
        //}



        //// 检查是否在攻击范围内
        //protected virtual bool IsInAttackRange()
        //{
        //    if (target == null) return false;
        //    return Vector2.Distance(transform.position, target.transform.position) <= attackRange;
        //}

        //// 尝试攻击
        //protected virtual void TryAttack()
        //{
        //    if (Time.time - lastAttackTime >= attackCooldown)
        //    {
        //        Attack();
        //    }
        //}

        //// 攻击逻辑
        //public virtual void Attack()
        //{
        //    isAttacking = true;
        //    animator.SetTrigger("Attack");
        //    lastAttackTime = Time.time;

        //    // 播放攻击特效
        //    if (attackEffect != null && attackPoint != null)
        //    {
        //        Instantiate(attackEffect, attackPoint.position, Quaternion.identity);
        //    }

        //    // 对目标造成伤害
        //    if (target != null && IsInAttackRange())
        //    {
        //        target.TakeDamage(attackPower);
        //    }

        //    Invoke(nameof(EndAttack), 0.5f);
        //}

        //// 结束攻击状态
        //protected virtual void EndAttack()
        //{
        //    isAttacking = false;
        //}

        //// 受击逻辑
        //public virtual void TakeDamage(int damage)
        //{
        //    currentHealth = Mathf.Max(0, currentHealth - damage);
        //    isTakingDamage = true;
        //    animator.SetTrigger("TakeDamage");

        //    // 播放受击特效
        //    if (hitEffect != null)
        //    {
        //        Instantiate(hitEffect, transform.position, Quaternion.identity);
        //    }

        //    // 击退效果
        //    if (target != null)
        //    {
        //        Vector2 knockbackDirection = ((Vector2)transform.position - (Vector2)target.transform.position).normalized;
        //        rb.velocity = knockbackDirection * 3f;
        //    }

        //    if (currentHealth <= 0)
        //    {
        //        Die();
        //    }
        //    else
        //    {
        //        Invoke(nameof(EndTakingDamage), 0.3f);
        //    }
        //}

        //// 结束受击状态
        //protected virtual void EndTakingDamage()
        //{
        //    isTakingDamage = false;
        //    rb.velocity = Vector2.zero;
        //}

        //// 死亡逻辑
        //public virtual void Die()
        //{
        //    isAlive = false;
        //    isMoving = false;
        //    animator.SetTrigger("Die");
        //    rb.velocity = Vector2.zero;

        //    // 播放死亡特效
        //    if (deathEffect != null)
        //    {
        //        Instantiate(deathEffect, transform.position, Quaternion.identity);
        //    }

        //    // 移除目标引用
        //    if (target != null)
        //    {
        //        Combatant targetCombatant = target;
        //        target = null;
        //        if (targetCombatant.target == this)
        //        {
        //            targetCombatant.target = null;
        //        }
        //    }

        //    Destroy(gameObject, 1f);
        //}

        //// 使用技能
        //public virtual void UseSkill(int skillIndex)
        //{
        //    if (!isAlive || isAttacking || isTakingDamage) return;

        //    animator.SetTrigger($"Skill{skillIndex}");
        //    lastAttackTime = Time.time; // 技能也会触发冷却

        //    // 技能特效和逻辑由子类实现
        //}

        //// 设置目标
        //public virtual void SetTarget(Combatant newTarget)
        //{
        //    target = newTarget;
        //}
    }
}