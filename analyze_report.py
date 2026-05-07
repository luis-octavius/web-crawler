#!/usr/bin/env python3
"""
Análise de Relatório do Web Crawler

Este script analisa o CSV gerado pelo crawler e fornece estatísticas úteis.

Uso:
    python3 analyze_report.py report.csv
"""

import pandas as pd
import sys
from pathlib import Path


def print_header(title):
    """Imprime um cabeçalho formatado"""
    print("\n" + "=" * 60)
    print(f"  {title}")
    print("=" * 60)


def analyze_general_stats(df):
    """Analisa estatísticas gerais"""
    print_header("📊 ESTATÍSTICAS GERAIS")
    
    print(f"\n  Total de páginas crawleadas: {len(df)}")
    print(f"  Páginas com H1: {df['H1'].notna().sum()} ({df['H1'].notna().sum() / len(df) * 100:.1f}%)")
    print(f"  Páginas sem H1: {df['H1'].isna().sum()} ({df['H1'].isna().sum() / len(df) * 100:.1f}%)")
    print(f"  Páginas com parágrafo: {df['FirstParagraph'].notna().sum()} ({df['FirstParagraph'].notna().sum() / len(df) * 100:.1f}%)")


def analyze_links(df):
    """Analisa links de saída"""
    print_header("🔗 ANÁLISE DE LINKS")
    
    # Conta links por página
    df['num_links'] = df['OutgoingLinks'].str.count(',').fillna(0) + 1
    df.loc[df['OutgoingLinks'].isna(), 'num_links'] = 0
    
    print(f"\n  Média de links por página: {df['num_links'].mean():.2f}")
    print(f"  Mediana de links: {df['num_links'].median():.0f}")
    print(f"  Máximo de links: {df['num_links'].max():.0f}")
    print(f"  Mínimo de links: {df['num_links'].min():.0f}")
    
    # Top páginas com mais links
    print("\n  📈 Top 5 páginas com mais links:")
    top_links = df.nlargest(5, 'num_links')[['URL', 'num_links']]
    for idx, row in top_links.iterrows():
        print(f"     {row['num_links']:.0f} links - {row['URL'][:60]}...")


def analyze_images(df):
    """Analisa imagens"""
    print_header("🖼️  ANÁLISE DE IMAGENS")
    
    # Conta imagens por página
    df['num_images'] = df['ImageURLs'].str.count(',').fillna(0) + 1
    df.loc[df['ImageURLs'].isna(), 'num_images'] = 0
    
    print(f"\n  Média de imagens por página: {df['num_images'].mean():.2f}")
    print(f"  Total de imagens encontradas: {df['num_images'].sum():.0f}")
    print(f"  Páginas com imagens: {(df['num_images'] > 0).sum()} ({(df['num_images'] > 0).sum() / len(df) * 100:.1f}%)")
    print(f"  Páginas sem imagens: {(df['num_images'] == 0).sum()} ({(df['num_images'] == 0).sum() / len(df) * 100:.1f}%)")
    
    # Top páginas com mais imagens
    if df['num_images'].max() > 0:
        print("\n  📸 Top 5 páginas com mais imagens:")
        top_images = df.nlargest(5, 'num_images')[['URL', 'num_images']]
        for idx, row in top_images.iterrows():
            print(f"     {row['num_images']:.0f} imagens - {row['URL'][:60]}...")


def analyze_content(df):
    """Analisa conteúdo textual"""
    print_header("📝 ANÁLISE DE CONTEÚDO")
    
    # Comprimento dos H1s
    df['h1_length'] = df['H1'].str.len()
    h1_with_content = df[df['H1'].notna()]
    
    if len(h1_with_content) > 0:
        print(f"\n  Comprimento médio do H1: {h1_with_content['h1_length'].mean():.1f} caracteres")
        print(f"  H1 mais curto: {h1_with_content['h1_length'].min():.0f} caracteres")
        print(f"  H1 mais longo: {h1_with_content['h1_length'].max():.0f} caracteres")
    
    # Comprimento dos parágrafos
    df['paragraph_length'] = df['FirstParagraph'].str.len()
    para_with_content = df[df['FirstParagraph'].notna()]
    
    if len(para_with_content) > 0:
        print(f"\n  Comprimento médio do 1º parágrafo: {para_with_content['paragraph_length'].mean():.1f} caracteres")
        print(f"  Parágrafo mais curto: {para_with_content['paragraph_length'].min():.0f} caracteres")
        print(f"  Parágrafo mais longo: {para_with_content['paragraph_length'].max():.0f} caracteres")


def analyze_seo_issues(df):
    """Identifica possíveis problemas de SEO"""
    print_header("⚠️  POSSÍVEIS PROBLEMAS DE SEO")
    
    issues = []
    
    # Páginas sem H1
    no_h1 = df[df['H1'].isna()]
    if len(no_h1) > 0:
        issues.append(f"  ❌ {len(no_h1)} páginas sem H1")
        print(f"\n  Páginas sem H1 ({len(no_h1)}):")
        for url in no_h1['URL'].head(5):
            print(f"     - {url}")
        if len(no_h1) > 5:
            print(f"     ... e mais {len(no_h1) - 5} páginas")
    
    # Páginas sem conteúdo
    no_content = df[df['FirstParagraph'].isna()]
    if len(no_content) > 0:
        issues.append(f"  ❌ {len(no_content)} páginas sem parágrafo")
        print(f"\n  Páginas sem primeiro parágrafo ({len(no_content)}):")
        for url in no_content['URL'].head(5):
            print(f"     - {url}")
        if len(no_content) > 5:
            print(f"     ... e mais {len(no_content) - 5} páginas")
    
    # H1s muito curtos (< 10 caracteres)
    short_h1 = df[(df['H1'].notna()) & (df['H1'].str.len() < 10)]
    if len(short_h1) > 0:
        issues.append(f"  ⚠️  {len(short_h1)} páginas com H1 muito curto (< 10 caracteres)")
        print(f"\n  Páginas com H1 muito curto ({len(short_h1)}):")
        for idx, row in short_h1.head(5).iterrows():
            print(f"     - '{row['H1']}' em {row['URL'][:50]}...")
    
    # H1s muito longos (> 70 caracteres)
    long_h1 = df[(df['H1'].notna()) & (df['H1'].str.len() > 70)]
    if len(long_h1) > 0:
        issues.append(f"  ⚠️  {len(long_h1)} páginas com H1 muito longo (> 70 caracteres)")
        print(f"\n  Páginas com H1 muito longo ({len(long_h1)}):")
        for idx, row in long_h1.head(3).iterrows():
            print(f"     - '{row['H1'][:60]}...' em {row['URL'][:40]}...")
    
    # Páginas com poucos links internos (< 3)
    df['num_links'] = df['OutgoingLinks'].str.count(',').fillna(0) + 1
    few_links = df[df['num_links'] < 3]
    if len(few_links) > 0:
        issues.append(f"  ⚠️  {len(few_links)} páginas com poucos links internos (< 3)")
    
    if not issues:
        print("\n  ✅ Nenhum problema de SEO detectado!")
    else:
        print(f"\n  Total de problemas encontrados: {len(issues)}")


def generate_summary(df):
    """Gera um resumo executivo"""
    print_header("📋 RESUMO EXECUTIVO")
    
    # Calcula score de qualidade (0-100)
    score = 0
    
    # Páginas com H1 (30 pontos)
    h1_score = (df['H1'].notna().sum() / len(df)) * 30
    score += h1_score
    
    # Páginas com conteúdo (30 pontos)
    content_score = (df['FirstParagraph'].notna().sum() / len(df)) * 30
    score += content_score
    
    # Média de links (20 pontos - ideal: 5-15 links)
    df['num_links'] = df['OutgoingLinks'].str.count(',').fillna(0) + 1
    avg_links = df['num_links'].mean()
    if 5 <= avg_links <= 15:
        links_score = 20
    elif avg_links < 5:
        links_score = (avg_links / 5) * 20
    else:
        links_score = max(0, 20 - (avg_links - 15))
    score += links_score
    
    # Páginas com imagens (20 pontos)
    df['num_images'] = df['ImageURLs'].str.count(',').fillna(0) + 1
    images_score = ((df['num_images'] > 0).sum() / len(df)) * 20
    score += images_score
    
    print(f"\n  Score de Qualidade: {score:.1f}/100")
    print(f"\n  Breakdown:")
    print(f"    - H1s: {h1_score:.1f}/30")
    print(f"    - Conteúdo: {content_score:.1f}/30")
    print(f"    - Links: {links_score:.1f}/20")
    print(f"    - Imagens: {images_score:.1f}/20")
    
    # Classificação
    if score >= 80:
        classification = "✅ EXCELENTE"
    elif score >= 60:
        classification = "👍 BOM"
    elif score >= 40:
        classification = "⚠️  REGULAR"
    else:
        classification = "❌ PRECISA MELHORAR"
    
    print(f"\n  Classificação: {classification}")


def main():
    """Função principal"""
    if len(sys.argv) < 2:
        print("Uso: python3 analyze_report.py report.csv")
        sys.exit(1)
    
    csv_file = sys.argv[1]
    
    # Verifica se o arquivo existe
    if not Path(csv_file).exists():
        print(f"❌ Erro: Arquivo '{csv_file}' não encontrado!")
        sys.exit(1)
    
    # Carrega o CSV
    try:
        df = pd.read_csv(csv_file)
    except Exception as e:
        print(f"❌ Erro ao ler CSV: {e}")
        sys.exit(1)
    
    # Verifica se o CSV tem dados
    if len(df) == 0:
        print("❌ Erro: CSV está vazio!")
        sys.exit(1)
    
    # Executa análises
    print("\n" + "🕷️  " * 20)
    print("  ANÁLISE DE RELATÓRIO DO WEB CRAWLER")
    print("🕷️  " * 20)
    
    analyze_general_stats(df)
    analyze_links(df)
    analyze_images(df)
    analyze_content(df)
    analyze_seo_issues(df)
    generate_summary(df)
    
    print("\n" + "=" * 60)
    print("  Análise concluída!")
    print("=" * 60 + "\n")


if __name__ == "__main__":
    main()
