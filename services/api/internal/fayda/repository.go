package fayda

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgRepository struct{ pool *pgxpool.Pool }

func NewPgRepository(pool *pgxpool.Pool) *PgRepository { return &PgRepository{pool: pool} }

func (r *PgRepository) Create(ctx context.Context, v *FaydaVerification) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO fayda_verifications (id, merchant_id, owner_id, kyc_profile_id, fin_hash, fin_last4, fan, partner_code, request_id, fayda_transaction_id, verification_method, otp_requested_at, consent_timestamp, consent_ip, status, front_doc_id, back_doc_id, selfie_doc_id, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18, now(), now())`,
		v.ID, v.MerchantID, v.OwnerID, v.KYCProfileID, v.FinHash, v.FinLast4, v.FAN, v.PartnerCode, v.RequestID, v.FaydaTransactionID, v.VerificationMethod, v.OTPRequestedAt, v.ConsentTimestamp, v.ConsentIP, v.Status, v.FrontDocID, v.BackDocID, v.SelfieDocID)
	return err
}
func (r *PgRepository) GetByRequestID(ctx context.Context, requestID string) (*FaydaVerification, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, owner_id, kyc_profile_id, fin_hash, fin_last4, fan, partner_code, request_id, fayda_transaction_id, verification_method, consent_timestamp, status, otp_verified, face_match_score FROM fayda_verifications WHERE request_id=$1`, requestID)
	var v FaydaVerification
	err := row.Scan(&v.ID, &v.MerchantID, &v.OwnerID, &v.KYCProfileID, &v.FinHash, &v.FinLast4, &v.FAN, &v.PartnerCode, &v.RequestID, &v.FaydaTransactionID, &v.VerificationMethod, &v.ConsentTimestamp, &v.Status, &v.OTPVerified, &v.FaceMatchScore)
	return &v, err
}
func (r *PgRepository) GetByOwner(ctx context.Context, ownerID string) ([]FaydaVerification, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, merchant_id, fin_last4, status, otp_verified, face_match_score FROM fayda_verifications WHERE owner_id=$1 ORDER BY created_at DESC`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []FaydaVerification
	for rows.Next() {
		var v FaydaVerification
		if err := rows.Scan(&v.ID, &v.MerchantID, &v.FinLast4, &v.Status, &v.OTPVerified, &v.FaceMatchScore); err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, nil
}
func (r *PgRepository) UpdateStatus(ctx context.Context, requestID string, status VerificationStatus, txID string, otpVerified bool) error {
	_, err := r.pool.Exec(ctx, `UPDATE fayda_verifications SET status=$1, fayda_transaction_id=$2, otp_verified=$3, updated_at=now() WHERE request_id=$4`, status, txID, otpVerified, requestID)
	return err
}
func (r *PgRepository) UpdateVerificationResult(ctx context.Context, requestID string, demoMatch bool, faceMatch bool, faceScore float64, encryptedRef string) error {
	_, err := r.pool.Exec(ctx, `UPDATE fayda_verifications SET demographics_match=$1, face_match=$2, face_match_score=$3, response_encrypted_ref=$4, verified_at=now(), updated_at=now() WHERE request_id=$5`, demoMatch, faceMatch, faceScore, encryptedRef, requestID)
	return err
}
